package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go-lb/pkg/logger"
)

type HealthChecker struct {
	Backends []string
	Status   map[int]bool // backend index -> healthy
	mu       sync.RWMutex
	interval time.Duration
	path     string
	log      logger.Logger
}

func NewHealthChecker(
	backends []string,
	interval time.Duration,
	path string,
	log logger.Logger,
) *HealthChecker {
	hc := &HealthChecker{
		Backends: backends,
		Status:   make(map[int]bool, len(backends)),
		interval: interval,
		path:     path,
		log:      log,
	}
	for i := range backends {
		hc.Status[i] = true // assume healthy at start
	}
	return hc
}

func (hc *HealthChecker) Start() {
	for i, backend := range hc.Backends {
		go hc.checkLoop(i, backend)
	}
}

func (hc *HealthChecker) checkLoop(idx int, backend string) {
	prevHealthy := true
	url := backend
	if !(len(url) > 4 && (url[:4] == "http")) {
		url = "http://" + url
	}
	url = url + hc.path
	for {
		time.Sleep(hc.interval)
		resp, err := http.Get(url)
		healthy := err == nil && resp.StatusCode < 500
		if resp != nil {
			resp.Body.Close()
		}
		hc.mu.Lock()
		hc.Status[idx] = healthy
		hc.mu.Unlock()
		if healthy != prevHealthy {
			if healthy {
				hc.log.Info(fmt.Sprintf("Backend %s is now healthy", backend))
			} else {
				hc.log.Warn(fmt.Sprintf("Backend %s is now UNHEALTHY", backend))
			}
			prevHealthy = healthy
		}
	}
}

func (hc *HealthChecker) HealthyIndexes() []int {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	idxs := make([]int, 0, len(hc.Status))
	for i, ok := range hc.Status {
		if ok {
			idxs = append(idxs, i)
		}
	}
	return idxs
}
