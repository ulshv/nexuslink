package tcp_client

import (
	"fmt"
	"net"

	"github.com/ulshv/nexuslink/internal/pb"
)

type ServerConnection struct {
	conn       net.Conn
	messagesCh chan *pb.TCPMessage
}

func NewServerConnection(host string, port int) (*ServerConnection, error) {
	// connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %v", err)
	}
	return &ServerConnection{
		conn: conn,
		// MessagesCh: make(chan *pb.TCPCommand),
	}, nil
}

// func RunClient(ctx context.Context, messagesCh chan<- *pb.TCPCommand, conn ServerConnection) {
// 	go ReadMessagesLoop(messagesCh, conn)
// }

func (c *ServerConnection) ListenAndHandle() {

}

func (c *ServerConnection) Close() error {
	return c.conn.Close()
}

// implements NetConnection interface
func (c *ServerConnection) Connection() net.Conn {
	return c.conn
}
