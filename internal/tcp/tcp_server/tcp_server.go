package tcp_server

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	// "github.com/ulshv/nexuslink/internal/pb"
)

type Server struct {
	host string
	port int
	// commandsCh chan *pb.TCPCommand
}

type ClientConnection struct {
	conn net.Conn
	// user *ClientUser
}

// type ClientUser struct {
// 	Username string
// 	Password string
// }

func NewServer(host string, port int) *Server {
	if host == "" {
		host = "0.0.0.0"
	}
	return &Server{
		host: host,
		port: port,
		// commandsCh: make(chan *pb.TCPCommand),
	}
}

func (s *Server) ListenAndHandle() {
	// Create a TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", s.host, s.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	slog.Info("server started and listening", "host", s.host, "port", s.port)
	for {
		// Accept a new client connection
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("failed to accept new client connection: %v", err)
			continue
		}
		slog.Info("new connection", "remote_addr", conn.RemoteAddr())
		// client := ClientConnection{
		// 	conn: conn,
		// }

		// slog.Info("new client", "client", client)
		// // server commands handler
		// go HandleServerSideCommands(context.Background(), server.messagesCh, client)
		// // read messages from the client (curent `conn`)
		// go ReadMessagesLoop(server.messagesCh, client)
	}
}

// // implements NetConnection interface
// func (s ClientConnection) Connection() net.Conn {
// 	return s.conn
// }
