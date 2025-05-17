package main

import (
	"github.com/goriiin/go-http-balancer/backend/configs"
	"github.com/goriiin/go-http-balancer/backend/internal/apps/global"
	del "github.com/goriiin/go-http-balancer/backend/internal/delivery"
	"github.com/goriiin/go-http-balancer/backend/internal/repository"
	"github.com/goriiin/go-http-balancer/pkg/logger"
	"github.com/goriiin/go-http-balancer/pkg/postgtresql"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const group = "group1"
const configFile = ".env"

func main() {
	err := godotenv.Load(configFile)
	if err != nil {
		log.Fatal("[ ERROR ] main.Load .env file")
	}

	conf, err := configs.GetAppConfig(group)
	if err != nil {
		log.Printf("[ ERROR ] main.GetAppConfig %s \n", err.Error())
	}

	appLogger := logger.InitLogger(conf.Logger)

	dbPool, err := postgtresql.NewPool(conf.DB)
	if err != nil {
		log.Fatalf("[ ERROR ] main.NewPool %s \n", err.Error())
	}
	defer dbPool.Close()

	repo := repository.NewDataRepository(dbPool, appLogger)

	delivery := del.NewDataDelivery(repo, appLogger)

	mux := http.NewServeMux()

	app := global.NewApp(delivery, mux, conf, appLogger)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := app.Run()

	select {
	case err := <-errChan:
		appLogger.Error("critical server failure", slog.String("error", err.Error()))
		app.Shutdown()

		os.Exit(1)

	case sig := <-sigChan:
		appLogger.Info("received shutdown signal", slog.String("signal", sig.String()))
		app.Shutdown()

		os.Exit(0)
	}
}
