package server

import (
	"fmt"
	"log"
	"net"
)

type TcpServer struct {
	host     string
	port     int
	listener net.Listener
}

func (s *TcpServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("server.go: Failed to start server: %v", err)
	}
	s.listener = listener
}

func (s *TcpServer) Stop() {
	s.listener.Close()
}

type NewTcpServerOpts struct {
	Host string
	Port int
}

func NewTcpServer(opts NewTcpServerOpts) *TcpServer {
	return &TcpServer{
		host: opts.Host,
		port: opts.Port,
	}
}
