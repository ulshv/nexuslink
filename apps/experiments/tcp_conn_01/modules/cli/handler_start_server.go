package cli

import (
	"context"
	"fmt"
	"net"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
	"github.com/ulshv/nexuslink/pkg/tcp_message"
)

func newHandlerStartServerCommand(lp *log_prompt.LogPrompt) cliHandlerFunc {
	logger := lp.NewLogger("handleStartServerCommand")
	return func(args []string) {
		handleStartServerCommand(logger, args)
	}
}

// args should contain []string{"<host>:<port>"} param
func handleStartServerCommand(logger log_prompt.Logger, args []string) {
	if len(args) == 0 {
		logger.Error("args are empty")
		return
	}
	host, port, err := parseHostPort(args[0])
	if err != nil {
		logger.Error("failed to parse host:port", "error", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		logger.Error("failed to listen", "host", host, "port", port, "error", err)
		return
	}
	logger.Info("server listening", "host", host, "port", port, "network", listener.Addr().Network())

	for {
		logger.Info("waiting for a new client connection")
		conn, err := listener.Accept()
		logger.Info("new client connection attempt", "remote_addr", conn.RemoteAddr())
		if err != nil {
			logger.Error("failed to accept new client connection", "error", err)
			continue
		}
		logger.Info("new client connection established", "remote_addr", conn.RemoteAddr())

		ctx, _ := context.WithCancel(context.Background())
		msgsChan := make(chan *tcp_message.TCPMessagePayload)

		go tcp_message.ReadTCPMessagesLoop(ctx, logger, msgsChan, conn)

		for msg := range msgsChan {
			logger.Info("server: received message from the client", "msg", msg)
		}
	}
}
