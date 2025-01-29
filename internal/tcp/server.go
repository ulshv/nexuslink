package tcp

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ulshv/nexuslink/internal/pb"
)

type Server struct {
	host       string
	port       string
	messagesCh chan *pb.TCPCommand
}

type ClientConnection struct {
	conn net.Conn
	user *ClientUser
}

type ClientUser struct {
	Username string
	Password string
}

type NewServerConfig struct {
	Host string
	Port string
}

func NewServer(config *NewServerConfig) *Server {
	return &Server{
		host:       config.Host,
		port:       config.Port,
		messagesCh: make(chan *pb.TCPCommand),
	}
}

func (server *Server) RunServer() {
	// Create a TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("[info]: server started, listening on %s:%s\n", server.host, server.port)

	for {
		// Accept a new client connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// logger.СlearCurrentLine()
		fmt.Println("[info]: new client connected")
		client := ClientConnection{
			conn: conn,
		}
		// server commands handler
		go HandleServerSideCommands(context.Background(), server.messagesCh, client)
		// read messages from the client (curent `conn`)
		go ReadMessagesLoop(server.messagesCh, client)
	}
}

// implements NetConnection interface
func (s ClientConnection) Connection() net.Conn {
	return s.conn
}
