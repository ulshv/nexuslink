package command

import (
	"context"
	"fmt"
	"sync"
)

func newCommandHandlers() CommandHandlers {
	return CommandHandlers{
		"help":    helloHandler,
		"start":   startHandler,
		"connect": connectHandler,
		"login":   loginHandler,
	}
}

func handleCommand(handlers CommandHandlers, command Command) {
	fmt.Println("[debug]: handleCommand: received command:", command)
	handler, ok := handlers[command.Command]
	if !ok {
		fmt.Println("[error]: handleCommand: command not found:", command.Command)
		return
	}
	handler(command.Args)
}

func CommandsWorker(ctx context.Context, wg *sync.WaitGroup, ch chan Command) {
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
