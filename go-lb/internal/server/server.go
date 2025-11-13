package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go-lb/internal/balancer"
	"go-lb/internal/config"
	"go-lb/pkg/logger"
)

func Run(cfg *config.Config, log logger.Logger) error {
	log.Info("The server runs with the following parameters:")
	log.Info(fmt.Sprintf("Mode: %s", cfg.Mode))
	log.Info(fmt.Sprintf("Front End IP address %s", cfg.FrontIP))
	log.Info(fmt.Sprintf("Backend addresses %s", cfg.BackendURLs))
	log.Info(fmt.Sprintf("Used algorithm: %s", cfg.Algorithm))
	log.Info(fmt.Sprintf("Health path %s", cfg.HealthPath))

	if cfg.Mode == "tcp" {
		interval, err := time.ParseDuration(cfg.HealthInterval)
		if err != nil {
			interval = 5 * time.Second
		}
		log.Info(fmt.Sprintf("Health interfval %s", interval))
		hc := NewTCPHealthChecker(cfg.BackendURLs, interval, log)
		hc.Start()
		var b balancer.Balancer
		switch cfg.Algorithm {
		case "weighted":
			b = balancer.NewWeighted(cfg.Weights)
		case "percentage":
			b = balancer.NewPercentage(cfg.Percentages)
		default:
			b = balancer.NewRoundRobin(len(cfg.BackendURLs))
		}
		ln, err := net.Listen("tcp", cfg.FrontIP)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf(
			"TCP proxy listening on %s, balancing to %v via %s",
			cfg.FrontIP,
			cfg.BackendURLs,
			cfg.Algorithm,
		))
		for {
			client, err := ln.Accept()
			if err != nil {
				log.Error(fmt.Sprintf("accept error: %v", err))
				continue
			}
			healthyIdxs := hc.HealthyIndexes()
			if len(healthyIdxs) == 0 {
				client.Close()
				continue
			}
			chosen := healthyIdxs[b.Next()%len(healthyIdxs)]
			backendAddr := cfg.BackendURLs[chosen]
			log.Debug(fmt.Sprintf("[TCP] Forwarding new connection to backend: %s", backendAddr))
			go func(client net.Conn, backendAddr string) {
				backend, err := net.Dial("tcp", backendAddr)
				if err != nil {
					client.Close()
					return
				}
				go io.Copy(backend, client)
				go io.Copy(client, backend)
			}(client, backendAddr)
		}
	} else {
		var b balancer.Balancer
		switch cfg.Algorithm {
		case "weighted":
			b = balancer.NewWeighted(cfg.Weights)
		case "percentage":
			b = balancer.NewPercentage(cfg.Percentages)
		default:
			b = balancer.NewRoundRobin(len(cfg.BackendURLs))
		}
		backends := make([]*url.URL, len(cfg.BackendURLs))
		for i, addr := range cfg.BackendURLs {
			u, err := url.Parse(addr)
			if err != nil {
				return err
			}
			// If no scheme, default to http
			if u.Scheme == "" {
				u, err = url.Parse("http://" + addr)
				if err != nil {
					return err
				}
			}
			backends[i] = u
		}
		interval, err := time.ParseDuration(cfg.HealthInterval)
		if err != nil {
			interval = 5 * time.Second
		}
		hc := NewHealthChecker(cfg.BackendURLs, interval, cfg.HealthPath, log)
		hc.Start()
		h := func(w http.ResponseWriter, r *http.Request) {
			log.Debug(fmt.Sprintf("[DEBUG] Handler called: %s %s", r.Method, r.URL.Path))

			healthyIdxs := hc.HealthyIndexes()
			if len(healthyIdxs) == 0 {
				log.Warn("No healthy backends available")
				w.WriteHeader(http.StatusServiceUnavailable)
				io.WriteString(w, "No healthy backends available")
				return
			}
			// Pick from healthy backends using the balancer
			chosen := healthyIdxs[b.Next()%len(healthyIdxs)]
			backendURL := backends[chosen]
			log.Debug(fmt.Sprintf("[HTTP] Forwarding request %s %s to backend: %s", r.Method, r.URL.Path, backendURL))
			proxy := httputil.NewSingleHostReverseProxy(backendURL)
			proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
				log.Error(fmt.Sprintf("proxy error: %v", err))
				w.WriteHeader(http.StatusBadGateway)
				io.WriteString(w, "Bad Gateway")
			}
			proxy.ServeHTTP(w, r)
		}
		log.Info(fmt.Sprintf(
			"Listening on %s, balancing to %v via %s",
			cfg.FrontIP,
			cfg.BackendURLs,
			cfg.Algorithm,
		))
		return http.ListenAndServe(cfg.FrontIP, http.HandlerFunc(h))
	}
}
