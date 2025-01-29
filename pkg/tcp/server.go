package tcp

import (
	"fmt"
	"log"
	"net"

	"github.com/ulshv/nexuslink/pkg/logger"
)

type Server struct {
	host string
	port string
}

type Connection struct {
	conn net.Conn
}

type NewServerConfig struct {
	Host string
	Port string
}

func NewServer(config *NewServerConfig) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) RunServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("[info]: server started, listening on %s:%s\n", server.host, server.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		logger.СlearCurrentLine()
		fmt.Printf("[info]: new client connected, Addr: %s\n", conn.LocalAddr().String())
		client := &Connection{
			conn: conn,
		}
		go client.handleRequestV2()
	}
}

func (client *Connection) handleRequestV2() {
	buf := make([]byte, 1024)
	for {
		n, err := client.conn.Read(buf)
		if err != nil {
			client.conn.Close()
			return
		}
		// Green console colour:    \x1b[32m
		// Reset console colour:    \x1b[0m
		fmt.Printf("\x1b[32m%s\x1b[0m%s", "[message]", ": ")
		fmt.Println(string(buf[:n]))
		fmt.Printf("[info]: tcp.client.message, length: %v\n", n)
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
