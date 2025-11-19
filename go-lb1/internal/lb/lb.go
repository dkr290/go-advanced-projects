package lb

import (
	"io"
	"log"
	"net"
)

type LoadBalancer struct {
	roundRobinCount int
	servers         []string
	n               int
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (l *LoadBalancer) Choose() string {
	idx := l.n % len(l.servers)
	l.n++
	return l.servers[idx]
}

func (l *LoadBalancer) Copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func (l *LoadBalancer) ServeProxy(client net.Conn, backendAddr string) {
	backend, err := net.Dial("tcp", backendAddr)
	if err != nil {
		client.Close()
		log.Printf("failed to dial %s: %s", backendAddr, err)
		return
	}

	go l.Copy(backend, client)
	go l.Copy(client, backend)
}
