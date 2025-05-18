package rate_limiter

import (
	"github.com/goriiin/go-http-balancer/balancer/configs"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func getClientID(r *http.Request, config configs.RateLimiterConfig, logger *slog.Logger) string {
	if config.ClientIDHeader != "" {
		apiKey := r.Header.Get(config.ClientIDHeader)
		if apiKey != "" {
			return apiKey
		}
	}

	var ip string
	if config.TrustXForwardedFor {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			ips := strings.Split(xff, ",")
			clientIP := strings.TrimSpace(ips[0])
			remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
			isTrusted := false

			if len(config.TrustedProxies) == 0 {
				isTrusted = true
			} else {
				parsedRemoteIP := net.ParseIP(remoteIP)

				if parsedRemoteIP != nil {
					for _, trustedCIDR := range config.TrustedProxies {
						_, ipNet, err := net.ParseCIDR(trustedCIDR)
						if err != nil {
							logger.Warn("invalid trusted_proxy CIDR in config", slog.String("cidr", trustedCIDR), slog.String("error", err.Error()))

							continue
						}
						if ipNet.Contains(parsedRemoteIP) {
							isTrusted = true

							break
						}
					}
				}
			}

			if isTrusted {
				ip = clientIP
			}
		}
	}

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
			logger.Debug("could not parse host from RemoteAddr, using full RemoteAddr", slog.String("remote_addr", r.RemoteAddr))
		} else {
			ip = host
		}
	}

	if ip == "" {
		logger.Warn("unable to determine client IP for rate limiting", slog.String("remote_addr", r.RemoteAddr), slog.String("x_forwarded_for", r.Header.Get("X-Forwarded-For")))

		return "unknown_client_ip"
	}

	return ip
}
