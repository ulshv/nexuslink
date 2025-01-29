package main

import (
	"context"
	"fmt"
	"sync"
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

	commandCh := make(chan Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	logHello()
	go commandsWorker(ctx, wg, commandCh)
	readCommandsLoop(commandCh)

	wg.Wait()
}
