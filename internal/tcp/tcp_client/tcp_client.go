package tcp_client

import (
	"fmt"
	"net"
	// "github.com/ulshv/nexuslink/internal/pb"
)

type ServerConnection struct {
	Conn net.Conn
	// commandsCh chan *pb.TCPCommand
}

type NewClientConfig struct {
	ServerHost string
	ServerPort string
}

func NewServerConnection(config NewClientConfig) (*ServerConnection, error) {
	// connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("[error]: failed to connect to server: %v", err)
	}
	return &ServerConnection{
		Conn: conn,
		// MessagesCh: make(chan *pb.TCPCommand),
	}, nil
}

// func RunClient(ctx context.Context, messagesCh chan<- *pb.TCPCommand, conn ServerConnection) {
// 	go ReadMessagesLoop(messagesCh, conn)
// }

func (c *ServerConnection) Close() error {
	return c.Conn.Close()
}

// implements NetConnection interface
func (c ServerConnection) Connection() net.Conn {
	return c.Conn
}
