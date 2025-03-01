package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
	"github.com/ulshv/nexuslink/pkg/tcp_message"
	"github.com/ulshv/nexuslink/pkg/tcp_message/pb"
)

// `prompt` usually have the following look:
// `command param1 param2 ...`
func HandlePrompt(lp *log_prompt.LogPrompt, prompt string) {
	logger := lp.NewLogger("prompt_handler")

	parts := strings.Split(prompt, " ")
	command := parts[0]
	params := parts[1:]

	switch command {
	case "server":
		handleServerCommand(lp, params)
	case "connect":
		handleConnectCommand(lp, params)
	case "help":
		logger.Log("Welcome to the NexusLink. Available commands:")
		logger.Log("	server <port> - start the server")
		logger.Log("  connect <host:port> - connect to the server")
		logger.Log("	help - show this message")
		logger.Log("	exit - exit the program")
	case "exit":
		logger.Log("Exiting...")
		os.Exit(0)
	default:
		logger.Log(fmt.Sprintf("Unknown command: %s", command))
	}
}

func handleServerCommand(lp *log_prompt.LogPrompt, params []string) {
	logger := lp.NewLogger("server_cmd_handler")

	if len(params) != 1 {
		logger.Log("server: wrong number of arguments")
		logger.Log("usage: server <port>")
		return
	}

	port := params[0]
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start server on port %s: %s", port, err))
		return
	}

	logger.Log(fmt.Sprintf("Server started on port %s", port))

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Error("Failed to accept connection", "error", err)
				continue
			}
			logger.Info("Accepted connection", "remote_addr", conn.RemoteAddr())
			handleConnection(lp, conn)
		}
	}()
}

func handleConnection(lp *log_prompt.LogPrompt, conn net.Conn) {
	ctx := context.Background()
	logger := lp.NewLogger("server_conn_handler")
	msgCh := make(chan *pb.TCPMessagePayload)

	go func() {
		defer conn.Close()
		go tcp_message.ReadTCPMessagesLoop(ctx, logger, msgCh, conn)
	}()

	go func() {
		for msg := range msgCh {
			logger.Info("Received message", "type", msg.Type, "data", string(msg.Data))
		}
	}()
}

func handleConnectCommand(lp *log_prompt.LogPrompt, params []string) {
	logger := lp.NewLogger("connect_cmd_handler")

	if len(params) != 1 {
		logger.Log("connect: wrong number of arguments")
		logger.Log("usage: connect <host:port>")
		return
	}

	port := params[0]
	conn, err := net.Dial("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error("Failed to connect to server", "error", err)
		return
	}

	payload := &pb.TCPMessagePayload{
		Type: "hello",
		Data: []byte("hello, world! what's up?"),
	}
	msg, _ := tcp_message.NewTCPMessage(lp.NewLogger("client"), payload)
	conn.Write(msg)
}
