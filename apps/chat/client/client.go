package main

import (
	"context"

	"github.com/ulshv/nexuslink/pkg/logprompt"
)

type ChatClient struct {
	lp *logprompt.LogPrompt
}

func NewChatClient(ctx context.Context) *ChatClient {
	lp := logprompt.NewLogPrompt(ctx, "> ")

	return &ChatClient{
		lp: lp,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewChatClient(ctx)

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
