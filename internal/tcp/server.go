package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/internal/pb"
	"google.golang.org/protobuf/proto"
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
		go client.handleRequestV3()
	}
}

func (client *Connection) handleRequestV3() {
	reader := bufio.NewReader(client.conn)
	for {
		// Read until we find the closing ":"
		prefix, err := reader.ReadBytes(':')
		if err != nil {
			client.conn.Close()
			return
		}

		// Verify the prefix starts with "protobuf("
		if !bytes.HasPrefix(prefix, []byte("protobuf(")) {
			fmt.Printf("[error]: invalid message prefix, expected 'protobuf('\n")
			return
		}

		// Parse the size from prefix "protobuf(%d):"
		end := bytes.Index(prefix, []byte(")"))
		if end == -1 || end != len(prefix)-2 { // -2 because of "):"
			fmt.Printf("[error]: invalid TCP message format! prefix: %s\n", prefix)
			return
		}

		size, err := strconv.Atoi(string(prefix[9:end])) // 9 is len("protobuf(")
		if err != nil {
			fmt.Printf("[error]: invalid size in prefix: %v\n", err)
			return
		}

		// Read exact number of bytes for the protobuf message
		data := make([]byte, size)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			fmt.Printf("[error]: failed to read message body: %v\n", err)
			return
		}

		command := &pb.TCPCommand{}
		err = proto.Unmarshal(data, command)
		if err != nil {
			fmt.Printf("[error]: failed to unmarshal message: %v\n", err)
			return
		}

		fmt.Printf(
			"[info]: received command `%s` with payload size of %d bytes.\n",
			command.Command,
			len(command.Payload),
		)

		// fmt.Printf("\x1b[32m%s\x1b[0m Protobuf message size: %d bytes\n", "[message]:", size)
		client.conn.Write([]byte("Message received.\n"))
	}
}

// func (client *Connection) handleRequestV2() {
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := client.conn.Read(buf)
// 		if err != nil {
// 			client.conn.Close()
// 			return
// 		}
// 		// Green console colour:    \x1b[32m
// 		// Reset console colour:    \x1b[0m
// 		fmt.Printf("\x1b[32m%s\x1b[0m%s", "[message]", ": ")
// 		fmt.Println(string(buf[:n]))
// 		fmt.Printf("[info]: tcp.client.message, length: %v\n", n)
// 		client.conn.Write([]byte("Message received.\n"))
// 	}
// }

// func (client *Connection) handleRequestV1() {
// 	reader := bufio.NewReader(client.conn)
// 	for {
// 		data, err := reader.ReadBytes('\n')
// 		if err != nil {
// 			client.conn.Close()
// 			return
// 		}
// 		// Green console colour:    \x1b[32m
// 		// Reset console colour:    \x1b[0m
// 		fmt.Printf("\x1b[32m%s\x1b[0m%s", "[message]", ": ")
// 		fmt.Println(string(data))
// 		// fmt.Printf("[info]: tcp.client.message, length: %v\n")
// 		client.conn.Write([]byte("Message received.\n"))
// 	}
// }
