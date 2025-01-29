package main

type Command struct {
	Command string
	Args    []string
}

type commandHandler func(args []string)

type commandHandlers map[string]commandHandler
