package tcp_message

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

func TestNewTCPMessage(t *testing.T) {
	t.Skip()
	t.Run("test message is correctly encoded", func(t *testing.T) {
		payload := &pb.TCPMessagePayload{
			Type: "hello",
			Data: []byte("hello, world! what's up?"),
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
		actualLength := string(msg[len(expectedPrefix):])
		if actualLength != string(payloadBytes) {
			t.Errorf("Expected message length to be %q but got %q", string(payloadBytes), actualLength)
		}
	})
}

func TestReadTCPMessagesLoop(t *testing.T) {
	tcpRW := &bytes.Buffer{}

	msgPayloads := []*pb.TCPMessagePayload{
		{Type: "hello", Data: []byte("")},
		{Type: "world", Data: []byte("")},
		{Type: "foo", Data: []byte("")},
		{Type: "bar", Data: []byte("")},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second) // to see the pause interval debug logs
		for _, payload := range msgPayloads {
			msg, err := NewTCPMessage(payload)
			if err != nil {
				t.Error(err)
			}
			tcpRW.Write(msg)
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond) // simple waiting for ReadTCPMessagesLoop to process the last msg
	}()

	msgPayloadsCh := make(chan *pb.TCPMessagePayload)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ReadTCPMessagesLoop(ctx, msgPayloadsCh, tcpRW)

	wg.Add(1)
	go func() {
		defer wg.Done()
		expectedPayloadIdx := 0
		for msg := range msgPayloadsCh {
			if msg.Type != msgPayloads[expectedPayloadIdx].Type {
				t.Errorf("Expected message type to be %q but got %q", msgPayloads[expectedPayloadIdx].Type, msg.Type)
			}
			expectedPayloadIdx++
		}
		if expectedPayloadIdx != len(msgPayloads) {
			t.Errorf("Expected %d messages but got %d", len(msgPayloads), expectedPayloadIdx)
		}
	}()

	wg.Wait()
}

// TODO: add test for partial data in the middle of the message, i.e. Write(msg[:len(msg)/2]), Write(msg[len(msg)/2:])
