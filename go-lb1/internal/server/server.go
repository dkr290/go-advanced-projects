// Package server is the main package to initialize and run lb
package server

import (
	"log"
	"net"
	"strings"
)

type ServerInf interface{}

type Server struct {
	ip      string
	port    string
	servers []string
}

func NewServer(ip string, port string, servers string) *Server {
	s := strings.Split(servers, ",")
	if len(s) == 1 && s[0] == "" {
		log.Fatalln("please specify backend servers with -backends")
	}
	return &Server{
		ip:      ip,
		port:    port,
		servers: s,
	}
}

func (s *Server) GetServers() []string {
	return s.servers
}

func (s *Server) Run() (net.Listener, error) {
	ln, err := net.Listen("tcp", s.ip+":"+s.port)
	log.Printf("listening on %s:%s, balancing %s", s.ip, s.port, s.servers)
	if err != nil {
		return nil, err
	}
	return ln, nil
}
