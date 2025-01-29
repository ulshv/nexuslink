package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/ulshv/nexuslink/internal/pb"
	"google.golang.org/protobuf/proto"
)

func SendMessage(conn NetConnection, command *pb.TCPCommand) error {
	// Marshal command to bytes
	data, err := proto.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}
	// Create prefix with data's length (in bytes)
	prefix := []byte(fmt.Sprintf("protobuf(%d):", len(data)))
	// Combine prefix + protobuf data
	message := append(prefix, data...)
	_, err = conn.Connection().Write(message)
	if err != nil {
		return fmt.Errorf("[error]: failed to send message: %w", err)
	}
	return nil
}

func ReadMessagesLoop(ch chan<- *pb.TCPCommand, conn NetConnection) {
	reader := bufio.NewReader(conn.Connection())
	for {
		// Read until we find the closing ":" part of the "protobuf(%d):" prefix
		prefix, err := reader.ReadBytes(':')
		if err != nil {
			conn.Connection().Close()
			return
		}
		// Verify the prefix starts with "protobuf("
		isPrefixValid := bytes.HasPrefix(prefix, []byte("protobuf("))
		if !isPrefixValid {
			fmt.Println("[error]: invalid message prefix, expected `protobuf(`")
			return
		}
		// Parse the size from prefix "protobuf(%d):"
		end := bytes.Index(prefix, []byte(")"))
		if end == -1 || end != len(prefix)-2 { // -2 because of "):"
			fmt.Printf("[error]: invalid TCP message format! prefix: %s\n", prefix)
			return
		}
		// Parse the size from prefix "protobuf(%d):"
		size, err := strconv.Atoi(string(prefix[9:end])) // 9 is len("protobuf(")
		if err != nil {
			fmt.Printf("[error]: invalid size in prefix: %v\n", err)
			return
		}
		// Read the binary data of the message from the reader
		data := make([]byte, size)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			fmt.Printf("[error]: failed to read message body: %v\n", err)
			return
		}
		// Try to unmarshal the data into a TCPCommand
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
		ch <- command
	}
}
