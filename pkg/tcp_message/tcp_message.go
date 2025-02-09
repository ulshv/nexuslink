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
	"time"

	"github.com/ulshv/nexuslink/internal/logger"
	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

type TCPMessage []byte

var newMsgLogger = logger.NewLogger("tcp_message/new_message")
var readMsgLogger = logger.NewLogger("tcp_message/read_messages")
var readMsgHeaderLogger = logger.NewLogger("tcp_message/read_messages_header")

const (
	messageHeaderPrefix = "_protobuf_("
	messageHeaderSuffix = "):"
	maxPayloadSize      = 1024 * 1024 // 1MB, adjust if needed
)

func NewTCPMessage(payload *pb.TCPMessagePayload) (TCPMessage, error) {
	newMsgLogger.Debug("NewTCPMessage", "message", payload)
	// Marshal payload to bytes
	data, err := proto.Marshal(payload)
	newMsgLogger.Debug("Marshalled payload to bytes", "bytes", len(data), "error", err)
	if err != nil {
		newMsgLogger.Error("Failed to marshal payload", "error", err)
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	// Create prefix with data's length (in bytes)
	header := []byte(fmt.Sprintf("%s%d%s", messageHeaderPrefix, len(data), messageHeaderSuffix))
	// Combine prefix + protobuf data
	msg := append(header, data...)
	newMsgLogger.Debug("Made TCPMessage, returning", "bytes", len(msg))
	return msg, nil
}

func ReadTCPMessagesLoop(
	ctx context.Context,
	ch chan<- *pb.TCPMessagePayload,
	r io.Reader,
) {
	readMsgLogger.Info("ReadTCPMessages")
	reader := bufio.NewReader(r)
	// messageHeadersCh := make(chan []byte)
	// go readTCPMessageHeadersLoop(ctx, messageHeadersCh, reader)
	readBytesEofCounter := 0
	lastReadBytesPauseInterval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			readMsgLogger.Info("context cancelled, exiting ReadTCPMessagesLoop")
			return
		// case messageHeader, ok := <-messageHeadersCh:
		// 	if !ok {
		// 		readMsgLogger.Info("messageHeadersCh closed, exiting ReadTCPMessagesLoop")
		// 		return
		// 	}
		default:
			readMsgLogger.Debug("waiting for the next ':' byte in the buffer...")
			// ReadBytes doesn't block if there's no data in the reader, it returns io.EOF err
			messageHeader, err := reader.ReadBytes(':')
			if err != nil {
				if err == io.EOF {
					readBytesEofCounter++
					pauseInterval := getPauseInterval(readBytesEofCounter)
					if pauseInterval != lastReadBytesPauseInterval {
						readMsgLogger.Debug("decreasing pause interval after EOF", "eofCounter", readBytesEofCounter, "pause", pauseInterval.String())
						lastReadBytesPauseInterval = pauseInterval
					}
					time.Sleep(pauseInterval)
					continue
				}
				readMsgLogger.Error("failed to read message prefix", "error", err)
				continue
			}
			readMsgLogger.Debug("recieved messageHeader", "messageHeader", string(messageHeader))
			isPrefixValid := bytes.HasPrefix(messageHeader, []byte(messageHeaderPrefix))
			if !isPrefixValid {
				// Trim the messageHeader to 20 bytes if it's longer, to prevent huge logs
				haderTrimLength := len(messageHeader)
				if haderTrimLength > 20 {
					haderTrimLength = 20
				}
				readMsgLogger.Error("invalid message prefix", "prefix", string(messageHeader[:haderTrimLength]))
				continue
			}
			payloadSizeEndIdx := bytes.Index(messageHeader, []byte(messageHeaderSuffix))
			if payloadSizeEndIdx == -1 {
				readMsgLogger.Error("invalid TCP message format, can't calculate payload size", "prefix", string(messageHeader[:20]))
				continue
			}
			payloadSizeStr := string(messageHeader[len(messageHeaderPrefix):payloadSizeEndIdx])
			payloadSize, err := strconv.Atoi(payloadSizeStr)
			readMsgLogger.Debug("calculated payload size", "payloadSize", payloadSize)
			if err != nil {
				readMsgLogger.Error("payload size is not a number", "payloadSizeStr", payloadSizeStr, "error", err)
				continue
			}
			if payloadSize > maxPayloadSize {
				readMsgLogger.Error("payload size is too big", "payloadSize", payloadSize, "maxPayloadSize", maxPayloadSize)
				continue
			}
			readMsgLogger.Debug("start extraction of TCP message payload...")
			binPayload := make([]byte, payloadSize)
			readPayloadEofCounter := 0
			lastReadPayloadPauseInterval := time.Duration(0)
			payloadBytesRead := 0
			var readPayloadErr error
			for {
				n, err := reader.Read(binPayload[payloadBytesRead:])
				if err != nil {
					if err == io.EOF {
						readPayloadEofCounter++
						pauseInterval := getPauseInterval(readPayloadEofCounter)
						if pauseInterval != lastReadPayloadPauseInterval {
							readMsgLogger.Debug("decreasing pause interval after EOF", "eofCounter", readPayloadEofCounter, "pause", pauseInterval.String())
							lastReadPayloadPauseInterval = pauseInterval
						}
						time.Sleep(pauseInterval)
						continue
					}
					readPayloadErr = err
					break // exit from the for loop on any other error, except EOF
				}
				payloadBytesRead += n
				if payloadBytesRead == payloadSize {
					break
				}
			}
			if readPayloadErr != nil {
				readMsgLogger.Error("failed to read full message payload", "error", err)
				continue
			}
			readMsgLogger.Debug("extracted TCP message payload")
			payload := &pb.TCPMessagePayload{}
			err = proto.Unmarshal(binPayload, payload)
			if err != nil {
				readMsgLogger.Error("failed to unmarshal TCPMessagePayload", "error", err)
				continue
			}
			readMsgLogger.Info("received TCPMessagePayload, writing to channel", "payloadType", payload.Type, "dataBytes", len(payload.Data))
			ch <- payload
		}
	}
}

// func readTCPMessageHeadersLoop(ctx context.Context, ch chan []byte, reader *bufio.Reader) {
// 	defer close(ch)

// 	eofCounter := 0

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			readMsgHeaderLogger.Info("waiting for the next ':' byte in the buffer...")
// 			messageHeader, err := reader.ReadBytes(':')
// 			if err != nil {
// 				if err == io.EOF {
// 					eofCounter++
// 					readMsgHeaderLogger.Info("reader.ReadBytes() - EOF")
// 					pauseInterval := getPauseInterval(eofCounter)
// 					readMsgHeaderLogger.Debug("pausing after EOF", "pause", pauseInterval.String())
// 					time.Sleep(pauseInterval)
// 					continue
// 				}
// 				readMsgHeaderLogger.Error("failed to read message prefix", "error", err)
// 				continue
// 			}
// 			readMsgHeaderLogger.Debug("recieved messageHeader", "messageHeader", string(messageHeader))
// 			ch <- messageHeader
// 		}
// 	}
// }

// Pause intervals. Increase from 10ms for the first 100ms,
// then 100ms for the next second
// then 500ms indefinetley
func getPauseInterval(eofCounter int) time.Duration {
	ms10threshold := 10
	ms100threshold := ms10threshold + 10

	if eofCounter < ms10threshold {
		return 10 * time.Millisecond
	} else if eofCounter < ms100threshold {
		return 100 * time.Millisecond
	}
	return 500 * time.Millisecond
}
