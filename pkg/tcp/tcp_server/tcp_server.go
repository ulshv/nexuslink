package tcp_server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/internal/pb"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	host     string
	port     int
	listener net.Listener
	logger   *slog.Logger
}

const (
	clientPingInterval = 1 * time.Second
)

func NewServer(host string, port int) *Server {
	if host == "" {
		host = "0.0.0.0"
	}
	return &Server{
		host:   host,
		port:   port,
		logger: logger.NewLogger("Server"),
	}
}

func (s *Server) ListenAndHandle() {
	s.logger.Info("Server starting listening", "host", s.host, "port", s.port)
	// Create a TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", s.host, s.port))
	if err != nil {
		s.logger.Error("Failed to listen", "error", err)
		log.Fatal(err)
	}
	defer listener.Close()
	s.logger.Info("Server listening", "host", s.host, "port", s.port)
	for {
		s.logger.Info("Waiting for a new client connection")
		// Accept a new client connection
		conn, err := listener.Accept()
		s.logger.Info("New client connection attempt", "remote_addr", conn.RemoteAddr())
		if err != nil {
			s.logger.Error("Failed to accept new client connection", "error", err)
			continue
		}
		s.logger.Info("New client connection established", "remote_addr", conn.RemoteAddr())
	}
}

func (s *Server) Close() error {
	s.logger.Info("Closing server listener")
	if s.listener != nil {
		return s.listener.Close()
	}
	s.logger.Info("Server listener already closed")
	return nil
}

func (s *Server) pingClient(ctx context.Context, conn net.Conn) {
	s.logger.Debug("Sending ping message to client", "remote_addr", conn.RemoteAddr())
	msg := &pb.TCPMessage{
		Type: "ping",
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.logger.Error("Failed to marshal ping message", "error", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		s.logger.Error("Failed to send ping message to client", "error", err)
	}
}

func (s *Server) pingClientLoop(ctx context.Context, conn net.Conn) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Ping client loop context done", "remote_addr", conn.RemoteAddr())
				return
			case <-time.After(clientPingInterval):
				s.pingClient(ctx, conn)
			}
		}
	}()
}
