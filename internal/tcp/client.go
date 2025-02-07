package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/ulshv/nexuslink/internal/pb"
)

type ServerConnection struct {
	Conn       net.Conn
	MessagesCh chan *pb.TCPMessage
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
		Conn:       conn,
		MessagesCh: make(chan *pb.TCPMessage),
	}, nil
}

func RunClient(ctx context.Context, messagesCh chan<- *pb.TCPMessage, conn ServerConnection) {
	go ReadMessagesLoop(messagesCh, conn)
}

func (c *ServerConnection) Close() error {
	return c.Conn.Close()
}

// implements NetConnection interface
func (c ServerConnection) Connection() net.Conn {
	return c.Conn
}
