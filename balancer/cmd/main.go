package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/goriiin/go-http-balancer/balancer/configs"
	"github.com/goriiin/go-http-balancer/balancer/internal/delivery/back"
	"github.com/goriiin/go-http-balancer/balancer/internal/delivery/balancer"
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"github.com/goriiin/go-http-balancer/balancer/internal/middleware/breaker"
	"github.com/goriiin/go-http-balancer/balancer/internal/middleware/rate_limiter"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/heap"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/least_conns"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/random"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/round_robin"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/store"
	"github.com/goriiin/go-http-balancer/balancer/internal/usecase/health_checker"
	"github.com/goriiin/go-http-balancer/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	envFile           = ".env"
	configFileName    = "balancer"
	configFileType    = "yaml"
	configPath        = "./configs/"
	balancerConfigKey = "balancer"
)

type pool interface {
	Next() *domain.Backend
	Add(be *domain.Backend)
	MarkHealthy(be *domain.Backend)
	MarkUnhealthy(be *domain.Backend)
	Done(be *domain.Backend)
	GetAll() []*domain.Backend
}

func main() {
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("[ WARN ] main.Load .env file: %v (continuing without it)\n", err)
	}

	balAppCfg := configs.DefaultBalancerCfg()

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configPath)

	setViperDefaults(balancerConfigKey, balAppCfg)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("[ WARN ] BreakerConfig file '%s.%s' not found in '%s'. Using default values.\n", configFileName, configFileType, configPath)
		} else {
			log.Fatalf("[ FATAL ] Failed to read config file '%s': %v\n", viper.ConfigFileUsed(), err)
		}
	}

	if err := viper.UnmarshalKey(balancerConfigKey, &balAppCfg, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)

		dc.ErrorUnused = true
	}); err != nil {
		log.Fatalf("[ FATAL ] Failed to unmarshal '%s' configuration section: %v. Check config structure and types.\n", balancerConfigKey, err)
	}

	appLogger := logger.InitLogger(balAppCfg.Logger)
	appLogger.Info("balancer configuration loaded", slog.Any("effective_config", balAppCfg))

	var p pool
	switch strings.ToLower(balAppCfg.Algorithm) {
	case "least-connections":
		heapImpl := &heap.Heap{}
		p = least_conns.NewPool(heapImpl)
		appLogger.Info("using 'least-connections' balancing algorithm")
	case "random":
		p = random.NewPool()
		appLogger.Info("using 'random' balancing algorithm")
	case "round-robin":
		fallthrough
	default:
		p = round_robin.NewPool()
		appLogger.Info("using 'round-robin' balancing algorithm (default)")
	}

	lb := balancer.NewBalancer(p, appLogger)

	if len(balAppCfg.Backends) == 0 {
		appLogger.Warn("no backends configured. Balancer will not be able to route requests.")
	}
	for i, backendAddrStr := range balAppCfg.Backends {
		parsedURL, err := url.Parse(backendAddrStr)
		if err != nil {
			appLogger.Error("failed to parse backend URL, skipping",
				slog.String("url", backendAddrStr),
				slog.String("error", err.Error()))

			continue
		}

		domainBe := &domain.Backend{
			URL:   parsedURL,
			Index: i,
		}

		p.Add(domainBe)

		backendProxyLogger := appLogger.With(
			slog.String("component", "backend_proxy"),
			slog.String("backend_url", parsedURL.String()),
		)

		rawProxyHandler := back.NewBackend(domainBe, backendProxyLogger)
		cb := breaker.NewBreaker(balAppCfg.CircuitBreaker)
		handlerWithCB := cb.Middleware(rawProxyHandler)
		lb.AddBackendHandler(domainBe.URL.String(), handlerWithCB)

		appLogger.Info("configured backend",
			slog.String("url", backendAddrStr),
			slog.Int("initial_index", i))
	}

	hc := health_checker.New(p, balAppCfg.HealthCheck,
		appLogger.With(
			slog.String("component", "health_checker")),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(balAppCfg.Backends) > 0 && balAppCfg.HealthCheck.Period > 0 {
		go hc.Start(ctx)
	} else {
		appLogger.Info("health checker will not start (no backends configured or health check period is zero)")
	}

	rateLimiterLogger := appLogger.With(
		slog.String("component", "rate_limiter_store"))
	rateLimiterStore := store.NewMemoryStore(
		balAppCfg.RateLimiter.DefaultCapacity,
		balAppCfg.RateLimiter.DefaultRefillRatePerSecond,
		balAppCfg.RateLimiter.CleanupInterval,
		balAppCfg.RateLimiter.InactiveThreshold,
		rateLimiterLogger,
	)
	rateLimitMiddleware := rate_limiter.Middleware(
		rateLimiterStore,
		balAppCfg.RateLimiter,
		appLogger.With(
			slog.String("component", "rate_limiter_middleware")))

	finalHandler := rateLimitMiddleware(http.HandlerFunc(lb.Handle))

	httpServer := &http.Server{
		Addr:         balAppCfg.ListenAddr,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		appLogger.Info("starting balancer HTTP server", slog.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("http server ListenAndServe error: %w", err)
		}

		close(errChan)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err, ok := <-errChan:
		if ok && err != nil {
			appLogger.Error("critical server failure", slog.String("error", err.Error()))
			os.Exit(1)
		}

		appLogger.Info("server shut down (errChan closed or nil error)")

	case sig := <-sigChan:
		appLogger.Info("received shutdown signal", slog.String("signal", sig.String()))
		appLogger.Info("initiating graceful shutdown...")

		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), balAppCfg.GracefulShutdownTimeout)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			appLogger.Error("server graceful shutdown failed", slog.String("error", err.Error()))
		} else {
			appLogger.Info("server gracefully shut down")
		}

		appLogger.Info("Stopping rate limiter store...")
		rateLimiterStore.Stop()

		appLogger.Info("balancer stopped")
		os.Exit(0)
	}
}

func setViperDefaults(prefix string, cfg configs.BalancerAppConfig) {
	viper.SetDefault(fmt.Sprintf("%s.listen_addr", prefix), cfg.ListenAddr)
	viper.SetDefault(fmt.Sprintf("%s.algorithm", prefix), cfg.Algorithm)
	viper.SetDefault(fmt.Sprintf("%s.graceful_shutdown_timeout", prefix), cfg.GracefulShutdownTimeout.String())
	viper.SetDefault(fmt.Sprintf("%s.backends", prefix), cfg.Backends)

	viper.SetDefault(fmt.Sprintf("%s.logger.level", prefix), cfg.Logger.Level)

	viper.SetDefault(fmt.Sprintf("%s.health_check.path", prefix), cfg.HealthCheck.Path)
	viper.SetDefault(fmt.Sprintf("%s.health_check.period", prefix), cfg.HealthCheck.Period.String())
	viper.SetDefault(fmt.Sprintf("%s.health_check.timeout", prefix), cfg.HealthCheck.Timeout.String())

	viper.SetDefault(fmt.Sprintf("%s.circuit_breaker.failure_threshold", prefix), cfg.CircuitBreaker.FailureThreshold)
	viper.SetDefault(fmt.Sprintf("%s.circuit_breaker.reset_timeout", prefix), cfg.CircuitBreaker.ResetTimeout.String())
	viper.SetDefault(fmt.Sprintf("%s.circuit_breaker.half_open_requests", prefix), cfg.CircuitBreaker.HalfOpenRequests)

	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.default_capacity", prefix), cfg.RateLimiter.DefaultCapacity)
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.default_refill_rate_per_second", prefix), cfg.RateLimiter.DefaultRefillRatePerSecond)
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.client_id_header", prefix), cfg.RateLimiter.ClientIDHeader)
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.trust_x_forwarded_for", prefix), cfg.RateLimiter.TrustXForwardedFor)
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.trusted_proxies", prefix), cfg.RateLimiter.TrustedProxies)
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.cleanup_interval", prefix), cfg.RateLimiter.CleanupInterval.String())
	viper.SetDefault(fmt.Sprintf("%s.rate_limiter.inactive_threshold", prefix), cfg.RateLimiter.InactiveThreshold.String())
}
