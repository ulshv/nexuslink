// TCPMessage is an abstraction for a protobuf-encoded message blobs.
// It is used to send and receive messages over a TCP connection.
//
// TCPMessage consists of Multiple parts:
// - TCPMessageHeader: `protobuf(%d):` prefix where `%d` is the length of the TCPMessagePayload in bytes
// - TCPMessagePayload: protobuf-encoded message with `len(bytes) == %d`
// - TCPMessagePayload.type - any string, which is used to identify the type of the message
// - TCPMessagePayload.data - []byte, which is used to transport the actual message's data
package tcp_message

import (
	"fmt"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

type TCPMessage []byte

var protoLogger = logger.NewLogger("proto_stream")

const (
	messageHeaderPrefix = "_protobuf_("
	messageHeaderSuffix = "):"
)

func NewTCPMessage(payload *pb.TCPMessagePayload) (TCPMessage, error) {
	protoLogger.Debug("NewTCPMessage", "message", payload)
	// Marshal payload to bytes
	data, err := proto.Marshal(payload)
	protoLogger.Debug("Marshalled payload to bytes", "bytes", len(data), "error", err)
	if err != nil {
		protoLogger.Error("Failed to marshal payload", "error", err)
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	// Create prefix with data's length (in bytes)
	header := []byte(fmt.Sprintf("%s%d%s", messageHeaderPrefix, len(data), messageHeaderSuffix))
	// Combine prefix + protobuf data
	msg := append(header, data...)
	protoLogger.Debug("Made TCPMessage, returning", "bytes", len(msg))
	return msg, nil
}
