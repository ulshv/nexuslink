package main

import (
	"context"

	"github.com/ulshv/nexuslink/pkg/logprompt"
	"github.com/ulshv/nexuslink/pkg/tcp/tcp_client"
)

type ChatClient struct {
	lp         *logprompt.LogPrompt
	serverConn *tcp_client.ServerConnection
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := newChatClient(ctx)

	client.lp.Log("Welcome to NexusLink Chat! Type 'help' to see available commands.")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case prompt := <-client.lp.Prompts():
				handlePrompt(client, prompt)
			}
		}
	}()

	client.lp.Start()
}

func newChatClient(ctx context.Context) *ChatClient {
	lp := logprompt.NewLogPrompt(ctx, "> ")

	return &ChatClient{
		lp: lp,
	}
}

func (c *ChatClient) setServerConn(conn *tcp_client.ServerConnection) {
	c.serverConn = conn
}

func (c *ChatClient) sendMessage(message []byte) {
	if c.serverConn == nil {
		c.lp.Log("No server connection. Type 'connect <host:port>' to connect to a server first.")
		return
	}
	// tcp.PublishMessage(c.serverConn, &pb.TCPMessage{
	// 	Type:    "hello",
	// 	Payload: message,
	// })
}
