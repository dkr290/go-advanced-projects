package server

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPHealthChecker struct {
	Backends []string
	Status   map[int]bool
	mu       sync.RWMutex
	interval time.Duration
}

func NewTCPHealthChecker(backends []string, interval time.Duration) *TCPHealthChecker {
	hc := &TCPHealthChecker{
		Backends: backends,
		Status:   make(map[int]bool, len(backends)),
		interval: interval,
	}
	for i := range backends {
		hc.Status[i] = true
	}
	return hc
}

func (hc *TCPHealthChecker) Start() {
	for i, backend := range hc.Backends {
		go hc.checkLoop(i, backend)
	}
}

func (hc *TCPHealthChecker) checkLoop(idx int, backend string) {
	prevHealthy := true
	for {
		time.Sleep(hc.interval)
		conn, err := net.DialTimeout("tcp", backend, 2*time.Second)
		healthy := err == nil
		if conn != nil {
			conn.Close()
		}
		hc.mu.Lock()
		hc.Status[idx] = healthy
		hc.mu.Unlock()
		if healthy != prevHealthy {
			if healthy {
				log.Printf("TCP backend %s is now healthy", backend)
			} else {
				log.Printf("TCP backend %s is now UNHEALTHY", backend)
			}
			prevHealthy = healthy
		}
	}
}

func (hc *TCPHealthChecker) HealthyIndexes() []int {
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
