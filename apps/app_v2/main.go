package main

import (
	"context"
	"sync"

	"github.com/ulshv/nexuslink/apps/app_v2/cli"
	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lp := log_prompt.NewLogPrompt(ctx, "> ")

	cliModule := cli.NewCliModule(lp)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go lp.Start()

	wg.Add(1)
	go func() {
		promptHandler := cliModule.NewCliPromptsHandler()

		for prompt := range lp.Prompts() {
			promptHandler(prompt)
		}
	}()

	wg.Wait()
}
