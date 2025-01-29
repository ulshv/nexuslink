package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ulshv/nexuslink/internal/cli_app"
)

func logHello() {
	msg := `Welcome to the tcp-chat!
Type 'help' for a list of commands.
`
	fmt.Println(msg)
}

func main() {
	// main cancellation context (TODO)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// channel for CLI commands
	commandCh := make(chan cli_app.Command)
	// cli instance for interactive logging and prompting
	cli := cli_app.NewCLI("> ")
	// wait group for goroutines
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// print welcome message
	logHello()
	// start CLI commands worker
	go cli_app.CommandsWorker(ctx, wg, commandCh, cli)
	// start CLI reader
	go cli.ReadCommandsLoopV2(commandCh)
	fmt.Println("test!")
	go func() {
		fmt.Println("go func!")
		counter := 0
		for {
			counter++
			cli.Log("[debug] test cli log, counter: %v", counter)
			time.Sleep(1 * time.Second)
		}
	}()

	// wait for all goroutines to finish
	wg.Wait()
}
