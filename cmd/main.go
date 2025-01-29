package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ulshv/nexuslink/internal/cli"
)

func logHello() {
	msg := `Welcome to the tcp-chat!
Type 'help' for a list of commands.
`
	fmt.Println(msg)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commandCh := make(chan cli.Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	logHello()
	go cli.CommandsWorker(ctx, wg, commandCh)
	cli.ReadCommandsLoop(commandCh)

	wg.Wait()
}
