package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ulshv/nexuslink/internal/command"
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

	commandCh := make(chan command.Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	logHello()
	go command.CommandsWorker(ctx, wg, commandCh)
	command.ReadCommandsLoop(commandCh)

	wg.Wait()
}
