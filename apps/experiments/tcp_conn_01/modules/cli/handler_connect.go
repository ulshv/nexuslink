package cli

import (
	"fmt"
	"net"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
	"github.com/ulshv/nexuslink/pkg/tcp_message"
)

func newHandlerConnectCommand(lp *log_prompt.LogPrompt) func(args []string) {
	logger := lp.NewLogger("handleConnectCommand")
	return func(args []string) {
		handleConnectCommand(logger, args)
	}
}

func handleConnectCommand(logger log_prompt.Logger, args []string) {
	if len(args) == 0 {
		logger.Error("handleConnectCommand: args are empty")
		return
	}
	host, port, err := parseHostPort(args[0])
	if err != nil {
		logger.Error("handleConnectCommand: failed to parse host:port", "error", err)
		return
	}
	logger.Info("handleConnectCommand: connecting to the server", "host", host, "port", port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		logger.Error("handleConnectCommand: failed to connect to the server", "host", host, "port", port, "error", err)
		return
	}
	logger.Info("client: connected to the server", "host", host, "port", port, "network", conn.LocalAddr().Network())

	msg, err := tcp_message.NewTCPMessage(logger, &tcp_message.TCPMessagePayload{
		Type: "hello",
		Data: []byte("Hello, server!"),
	})
	if err != nil {
		logger.Error("client: failed to create new TCP message", "error", err)
		return
	}
	conn.Write(msg)
	logger.Info("client: sent message to the server")
}
