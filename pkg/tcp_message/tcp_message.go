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

	"github.com/ulshv/nexuslink/pkg/logs"
	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
	"google.golang.org/protobuf/proto"
)

type TCPMessage []byte
type TCPMessagePayload = pb.TCPMessagePayload

const (
	messageHeaderPrefix = "_protobuf_("
	messageHeaderSuffix = "):"
	maxPayloadSize      = 1024 * 1024 // 1MB, adjust if needed
)

func NewTCPMessage(logger logs.Logger, payload *pb.TCPMessagePayload) (TCPMessage, error) {
	logger.Debug("NewTCPMessage", "message", payload)
	data, err := proto.Marshal(payload)
	logger.Debug("Marshalled payload to bytes", "bytes", len(data), "error", err)
	if err != nil {
		logger.Error("Failed to marshal payload", "error", err)
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	header := []byte(fmt.Sprintf("%s%d%s", messageHeaderPrefix, len(data), messageHeaderSuffix))
	msg := append(header, data...)
	logger.Debug("Made TCPMessage, returning", "bytes", len(msg))
	return msg, nil
}

func ReadTCPMessagesLoop(
	ctx context.Context,
	logger logs.Logger,
	ch chan<- *pb.TCPMessagePayload,
	r io.Reader,
) {
	logger.Info("ReadTCPMessages")
	reader := bufio.NewReader(r)
	readBytesEofCounter := 0
	lastReadBytesPauseInterval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			logger.Info("context cancelled, closing channel and exiting ReadTCPMessagesLoop")
			close(ch)
			return
		default:
			// logger.Debug("waiting for the next ':' byte in the buffer...")

			// ReadBytes doesn't block if there's no data in the reader, it returns io.EOF err
			messageHeader, err := reader.ReadBytes(':')
			if err != nil {
				if err == io.EOF {
					readBytesEofCounter++
					pauseInterval := getPauseInterval(readBytesEofCounter)
					if pauseInterval != lastReadBytesPauseInterval {
						logger.Debug("decreasing pause interval after EOF", "eofCounter", readBytesEofCounter, "pause", pauseInterval.String())
						lastReadBytesPauseInterval = pauseInterval
					}
					time.Sleep(pauseInterval)
					continue
				}
				logger.Error("failed to read message prefix", "error", err)
				continue
			}
			logger.Debug("recieved messageHeader", "messageHeader", string(messageHeader))
			isPrefixValid := bytes.HasPrefix(messageHeader, []byte(messageHeaderPrefix))
			if !isPrefixValid {
				// Trim the messageHeader to 20 bytes if it's longer, to prevent huge logs
				haderTrimLength := len(messageHeader)
				if haderTrimLength > 20 {
					haderTrimLength = 20
				}
				logger.Error("invalid message prefix", "prefix", string(messageHeader[:haderTrimLength]))
				continue
			}
			payloadSizeEndIdx := bytes.Index(messageHeader, []byte(messageHeaderSuffix))
			if payloadSizeEndIdx == -1 {
				logger.Error("invalid TCP message format, can't calculate payload size", "prefix", string(messageHeader[:20]))
				continue
			}
			payloadSizeStr := string(messageHeader[len(messageHeaderPrefix):payloadSizeEndIdx])
			payloadSize, err := strconv.Atoi(payloadSizeStr)
			logger.Debug("calculated payload size", "payloadSize", payloadSize)
			if err != nil {
				logger.Error("payload size is not a number", "payloadSizeStr", payloadSizeStr, "error", err)
				continue
			}
			if payloadSize > maxPayloadSize {
				logger.Error("payload size is too big", "payloadSize", payloadSize, "maxPayloadSize", maxPayloadSize)
				continue
			}
			logger.Debug("start extraction of TCP message payload...")
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
							logger.Debug("decreasing pause interval after EOF", "eofCounter", readPayloadEofCounter, "pause", pauseInterval.String())
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
				logger.Error("failed to read full message payload", "error", err)
				continue
			}
			logger.Debug("extracted TCP message payload")
			payload := &pb.TCPMessagePayload{}
			err = proto.Unmarshal(binPayload, payload)
			if err != nil {
				logger.Error("failed to unmarshal TCPMessagePayload", "error", err)
				continue
			}
			logger.Info("received TCPMessagePayload, writing to channel", "payloadType", payload.Type, "dataBytes", len(payload.Data))
			ch <- payload
		}
	}
}

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
