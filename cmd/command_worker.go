package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ulshv/nexuslink/pkg/server"
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
- start <port>: Start a server
- connect <host:port>: Connect to a server
- disconnect <uuid>: Disconnect from the server
- list: List all connected servers
- info <uuid>: Show info about a server
- send <uuid> <message...>: Send a message to the server
- exit: Exit the program`)
		},
		"start": func(args []string) {
			if len(args) != 1 {
				fmt.Println("[error]: start: invalid number of arguments (need <port>)")
				return
			}
			fmt.Println("[info]: starting TCP server")
			server := server.New(&server.Config{
				Host: "0.0.0.0",
				Port: args[0],
			})
			go server.Run()
		},
	}
}
