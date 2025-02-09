// TCPMessage is an abstraction for a protobuf-encoded message blobs.
// It is used to send and receive messages over a TCP connection.
//
// TCPMessage consists of Multiple parts:
// - TCPMessageHeader: `_protobuf_(%d):` prefix
// - TCPMessagePayloadSize: `%d` from the TCPMessageHeader, it's the length of the TCPMessagePayload in bytes
// - TCPMessagePayload: protobuf-encoded message with `len(bytes) == %d from the TCPMessageHeader`
// - TCPMessagePayload.type - any string, which is used to identify the type of the message
// - TCPMessagePayload.data - []byte, which is used to transport the actual message's data
package tcp_message

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

type TCPMessage []byte

var tcpLogger = logger.NewLogger("tcp_message")

const (
	messageHeaderPrefix = "_protobuf_("
	messageHeaderSuffix = "):"
	maxPayloadSize      = 1024 * 1024 // 1MB, adjust if needed
)

func NewTCPMessage(payload *pb.TCPMessagePayload) (TCPMessage, error) {
	tcpLogger.Debug("NewTCPMessage", "message", payload)
	// Marshal payload to bytes
	data, err := proto.Marshal(payload)
	tcpLogger.Debug("Marshalled payload to bytes", "bytes", len(data), "error", err)
	if err != nil {
		tcpLogger.Error("Failed to marshal payload", "error", err)
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	// Create prefix with data's length (in bytes)
	header := []byte(fmt.Sprintf("%s%d%s", messageHeaderPrefix, len(data), messageHeaderSuffix))
	// Combine prefix + protobuf data
	msg := append(header, data...)
	tcpLogger.Debug("Made TCPMessage, returning", "bytes", len(msg))
	return msg, nil
}

func ReadTCPMessagesLoop(
	ctx context.Context,
	ch chan<- *pb.TCPMessagePayload,
	r io.Reader,
) {
	tcpLogger.Info("ReadTCPMessages")
	reader := bufio.NewReader(r)
	for {
		// loop:
		select {
		case <-ctx.Done():
			tcpLogger.Info("context cancelled, exiting ReadTCPMessagesLoop")
			return
		default:
			tcpLogger.Info("waiting for the next ':' byte in the buffer...")
			messageHeader, err := reader.ReadBytes(':') // read until ":" byte
			if err != nil {
				if err == io.EOF {
					tcpLogger.Info("reader.ReadBytes() - EOF")
					// break loop
					return
				}
				tcpLogger.Error("failed to read message prefix", "error", err)
				continue
			}
			tcpLogger.Debug("recieved messageHeader", "messageHeader", string(messageHeader))
			isPrefixValid := bytes.HasPrefix(messageHeader, []byte(messageHeaderPrefix))
			if !isPrefixValid {
				tcpLogger.Error("invalid message prefix", "prefix", string(messageHeader[:20]))
				continue
			}
			payloadSizeEndIdx := bytes.Index(messageHeader, []byte(messageHeaderSuffix))
			if payloadSizeEndIdx == -1 {
				tcpLogger.Error("invalid TCP message format, can't calculate payload size", "prefix", string(messageHeader[:20]))
				continue
			}
			payloadSizeStr := string(messageHeader[len(messageHeaderPrefix):payloadSizeEndIdx])
			payloadSize, err := strconv.Atoi(payloadSizeStr)
			tcpLogger.Debug("calculated payload size", "payloadSize", payloadSize)
			if err != nil {
				tcpLogger.Error("payload size is not a number", "payloadSizeStr", payloadSizeStr, "error", err)
				continue
			}
			if payloadSize > maxPayloadSize {
				tcpLogger.Error("payload size is too big", "payloadSize", payloadSize, "maxPayloadSize", maxPayloadSize)
				continue
			}
			tcpLogger.Debug("start extraction of TCP message payload...")
			binPayload := make([]byte, payloadSize)
			_, err = io.ReadFull(reader, binPayload)
			if err != nil {
				tcpLogger.Error("failed to read message payload", "error", err)
				continue
			}
			tcpLogger.Debug("extracted TCP message payload")
			payload := &pb.TCPMessagePayload{}
			err = proto.Unmarshal(binPayload, payload)
			if err != nil {
				tcpLogger.Error("failed to unmarshal TCPMessagePayload", "error", err)
				continue
			}
			tcpLogger.Info("received TCPMessagePayload, writing to channel", "payloadType", payload.Type, "dataBytes", len(payload.Data))
			ch <- payload
		}
	}
}
