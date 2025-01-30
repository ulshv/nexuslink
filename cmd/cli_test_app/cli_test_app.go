package main

import (
	"sync"
	"time"

	"github.com/ulshv/nexuslink/internal/cli_app"
)

func main() {
	// ch := make(chan cli_app.Command)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		counter := 0
		for {
			cli_app.LogV2("counter: %v", counter)
			counter++
			time.Sleep(1 * time.Second)
		}
	}()

	go cli_app.ReadCommandsLoopV4()

	wg.Wait()
}
