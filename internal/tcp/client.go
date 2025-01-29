package tcp

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

type NewClientConfig struct {
	Address string
	Port    string
}

func NewClient(config NewClientConfig) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.Address, config.Port))
	if err != nil {
		return nil, fmt.Errorf("[error]: failed to connect to server: %v", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) SendMessage(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("[error]: failed to send message: %v", err)
	}
	return nil
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
