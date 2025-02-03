package main

import "strings"

func handlePrompt(client *ChatClient, prompt string) {
	// trim whitespaces from the promtp (i.e. newlines, spaces in-beteen)
	cleanPrompt := strings.TrimSpace(prompt)
	// extract params from the prompt
	params := strings.Split(cleanPrompt, " ")
	// make sure that params have some strings
	if len(params) == 0 {
		client.lp.Log("[debug] empty prompt: %s", prompt)
		return
	}
	// extract command and params
	command := params[0]
	_ = params[1:]
	// handle commands
	switch command {
	case "help":
		handleHelp(client)
	default:
		handleUnknownCommand(client, command)
	}
}
