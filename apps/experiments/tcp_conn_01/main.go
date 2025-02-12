package main

import (
	"context"
	"sync"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

func main() {
	wg := sync.WaitGroup{}
	lp := log_prompt.NewLogPrompt(context.Background(), "> ")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for prompt := range lp.Prompts() {
			handleCliCommands(lp, prompt)
		}
	}()

	lp.Start()
	wg.Wait()
}
