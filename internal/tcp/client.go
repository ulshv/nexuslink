package tcp

import (
	"fmt"
	"net"

	"github.com/ulshv/nexuslink/internal/pb"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn net.Conn
}

type NewClientConfig struct {
	ServerHost string
	ServerPort string
}

func NewClient(config NewClientConfig) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("[error]: failed to connect to server: %v", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) ReceiveMessage() (string, error) {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("[error]: failed to receive message: %v", err)
	}
	return string(buf[:n]), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendMessageV2(command *pb.TCPCommand) error {
	// Marshal command to bytes
	data, err := proto.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %v", err)
	}

	// Create prefix with length
	prefix := []byte(fmt.Sprintf("protobuf(%d):", len(data)))

	// Combine prefix + protobuf data
	message := append(prefix, data...)

	_, err = c.conn.Write(message)
	if err != nil {
		return fmt.Errorf("[error]: failed to send message: %v", err)
	}
	return nil
}

// func (c *Client) SendMessageV1(command proto.TCPCommand) error {
// 	payload := append(command, []byte())
// 	_, err := c.conn.Write()
// 	if err != nil {
// 		return fmt.Errorf("[error]: failed to send message: %v", err)
// 	}
// 	return nil
// }
