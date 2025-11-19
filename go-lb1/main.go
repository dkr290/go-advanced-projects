package main

import (
	"log"

	"go-lb1/internal/lb"
	"go-lb1/internal/server"
)

func main() {
	server := server.NewServer("0.0.0.0", "8080", "10.167.10.146:26257,10.167.10.146:26257")
	ln, err := server.Run()
	if err != nil {
		log.Fatal(err)
	}

	lb := lb.NewLoadBalancer(server.GetServers())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		go lb.ServeProxy(conn, lb.Choose())
	}
}
