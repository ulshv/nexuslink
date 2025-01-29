package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ulshv/nexuslink/internal/cli_commands"
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

	commandCh := make(chan cli_commands.Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	logHello()
	go cli_commands.CommandsWorker(ctx, wg, commandCh)
	cli_commands.ReadCommandsLoop(commandCh)

	wg.Wait()
}
