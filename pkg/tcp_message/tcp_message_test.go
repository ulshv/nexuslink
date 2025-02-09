package tcp_message

import (
	"strconv"
	"testing"

	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

func TestNewTCPMessage(t *testing.T) {
	t.Run("test message is correctly encoded", func(t *testing.T) {
		payload := &pb.TCPMessagePayload{
			Type:    "hello",
			Payload: []byte("hello, world! what's up?"),
		}
		msg, err := NewTCPMessage(payload)
		if err != nil {
			t.Error(err)
		}
		payloadBytes, err := proto.Marshal(payload)
		if err != nil {
			t.Error(err)
		}
		expectedPrefix := messageHeaderPrefix + strconv.Itoa(len(payloadBytes)) + messageHeaderSuffix
		actualPrefix := string(msg[:len(expectedPrefix)])
		if actualPrefix != expectedPrefix {
			t.Errorf("Expected message prefix to be %q but got %q", expectedPrefix, actualPrefix)
		}
		actualPayload := string(msg[len(expectedPrefix):])
		if actualPayload != string(payloadBytes) {
			t.Errorf("Expected message payload to be %q but got %q", string(payloadBytes), actualPayload)
		}
	})
}
