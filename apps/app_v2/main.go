package main

import (
	"context"
	"sync"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lp := log_prompt.NewLogPrompt(ctx, "> ")
	mainLogger := lp.NewLogger("main")

	mainLogger.Info("Welcome to the NexusLink!")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go lp.Start()

	// cliModule := cli.NewCliModule(lp)

	for prompt := range lp.Prompts() {
		mainLogger.Info("received new", "prompt", prompt)
		// mainLogger.Log(fmt.Sprintf(`"%s", you say?`, prompt))
	}

	wg.Wait()
}
