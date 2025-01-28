package main

import (
	"context"
	"fmt"
	"sync"
)

type commandHandler func(args []string)
type commandHandlers map[string]commandHandler

func handleCommand(handlers commandHandlers, command Command) {
	fmt.Println("[debug]: handleCommand: received command:", command)
	handler, ok := handlers[command.Command]
	if !ok {
		fmt.Println("[error]: handleCommand: command not found:", command.Command)
		return
	}
	handler(command.Args)
}

func commandsWorker(ctx context.Context, wg *sync.WaitGroup, ch chan Command) {
	handlers := newCommandHandlers()
	defer wg.Done()

	for {
		select {
		case command := <-ch:
			handleCommand(handlers, command)
		case <-ctx.Done():
			fmt.Println("commandsWorker: received ctx cancellation, stopping the worker")
			return
		}
	}
}

func newCommandHandlers() commandHandlers {
	return commandHandlers{
		"help": func(args []string) {
			fmt.Println(`Available commands:
- help: Show this help message
- connect: Connect to a server
- disconnect: Disconnect from the server
- list: List all connected servers
- send: Send a message to the server
- exit: Exit the program`)
		},
	}
}
