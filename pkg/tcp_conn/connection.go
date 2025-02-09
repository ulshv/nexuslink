package tcp_conn

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/internal/pb"
	"google.golang.org/protobuf/proto"
)

type connection struct {
	conn   net.Conn
	logger *slog.Logger
}

func newConnection(conn net.Conn) *connection {
	return &connection{
		conn:   conn,
		logger: logger.NewLogger(fmt.Sprintf("Connection %s", conn.RemoteAddr())),
	}
}

func (c *connection) PublishMessage(message *pb.TCPMessage) error {
	c.logger.Info(
		"Publishing message",
		"message", message,
		"remote_addr", c.conn.RemoteAddr(),
	)
	// Marshal command to bytes
	data, err := proto.Marshal(message)
	if err != nil {
		c.logger.Error("Failed to marshal command", "error", err)
		return fmt.Errorf("failed to marshal command: %w", err)
	}
	// Create prefix with data's length (in bytes)
	prefix := []byte(fmt.Sprintf("protobuf(%d):", len(data)))
	// Combine prefix + protobuf data
	msg := append(prefix, data...)
	_, err = c.conn.Write(msg)
	if err != nil {
		return fmt.Errorf("[error]: failed to send message: %w", err)
	}
	return nil
}

func (c *connection) ReadMessagesLoop(ch chan<- *pb.TCPMessage) {
	reader := bufio.NewReader(c.conn)
	for {
		// Read until we find the closing ":" part of the "protobuf(%d):" prefix
		prefix, err := reader.ReadBytes(':')
		if err != nil {
			c.conn.Close()
			continue
		}
		// Verify the prefix starts with "protobuf("
		isPrefixValid := bytes.HasPrefix(prefix, []byte("protobuf("))
		if !isPrefixValid {
			fmt.Println("[error]: invalid message prefix, expected `protobuf(`")
			continue
		}
		// Parse the size from prefix "protobuf(%d):"
		end := bytes.Index(prefix, []byte(")"))
		if end == -1 || end != len(prefix)-2 { // -2 because of "):"
			fmt.Printf("[error]: invalid TCP message format! prefix: %s\n", prefix)
			continue
		}
		// Parse the size from prefix "protobuf(%d):"
		size, err := strconv.Atoi(string(prefix[9:end])) // 9 is len("protobuf(")
		if err != nil {
			fmt.Printf("[error]: invalid size in prefix: %v\n", err)
			continue
		}
		// Read the binary data of the message from the reader
		data := make([]byte, size)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			fmt.Printf("[error]: failed to read message body: %v\n", err)
			continue
		}
		// Try to unmarshal the data into a TCPMessage
		message := &pb.TCPMessage{}
		err = proto.Unmarshal(data, message)
		if err != nil {
			fmt.Printf("[error]: failed to unmarshal message: %v\n", err)
			continue
		}
		fmt.Printf(
			"[info]: received message `%s` with payload size of %d bytes.\n",
			message.Type,
			len(message.Payload),
		)
		ch <- message
	}
}
