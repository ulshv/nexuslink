package main

import (
	"context"
	"sync"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

func main() {
	appCtx := context.Background()
	lp := log_prompt.NewLogPrompt(appCtx, "> ")
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		lp.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-appCtx.Done():
				return
			case prompt := <-lp.Prompts():
				HandlePrompt(lp, prompt)
			}
		}
	}()

	wg.Wait()
}
