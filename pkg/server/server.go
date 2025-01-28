package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	host string
	port string
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host string
	Port string
}

func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequestV2()
	}
}

func (client *Client) handleRequestV2() {
	buf := make([]byte, 1024)
	for {
		n, err := client.conn.Read(buf)
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("tcp.client.message, length: %v, message: %v\n", n, string(buf[:n]))
		client.conn.Write([]byte("Message received.\n"))
	}
}

// func (client *Client) handleRequest() {
// 	reader := bufio.NewReader(client.conn)
// 	for {
// 		message, err := reader.ReadString('\n')
// 		if err != nil {
// 			client.conn.Close()
// 			return
// 		}
// 		fmt.Printf("Message incoming: %s", string(message))
// 		client.conn.Write([]byte("Message received.\n"))
// 	}
// }
