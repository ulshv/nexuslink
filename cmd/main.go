package main

import (
	"context"
	"fmt"
	"sync"
)

func initialize(port int) {
	fmt.Println("Welcome to the NexusLink!")
	fmt.Printf("Server is running on port %v\n", port)
	fmt.Printf("Type 'help' for a list of commands.\n\n")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commandCh := make(chan Command)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	initialize(8080)
	go commandsWorker(ctx, wg, commandCh)
	readCommandsLoop(commandCh)

	wg.Wait()
}
