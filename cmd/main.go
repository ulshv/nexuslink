package main

import (
	"context"
	"fmt"
	"sync"
)

func logHello(port int) {
	fmt.Printf(`Welcome to the NexusLink!
Server is running on port %v
Type 'help' for a list of commands.

`, port)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commandCh := make(chan Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	logHello(8080)
	go commandsWorker(ctx, wg, commandCh)
	readCommandsLoop(commandCh)

	wg.Wait()
}
