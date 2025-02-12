package cli

import (
	"strings"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

// [prompt] arg is a string with with a command and params.
// For example:
//
//	   "serve <host>:<port>"
//		 "connect <host>:<port>"
//		 "dm @username mesage text here"
func NewCliPromptsHandler(lp *log_prompt.LogPrompt) func(prompt string) {
	logger := lp.NewLogger("handlerCliCommands")

	handleStartServerCommand := newHandlerStartServerCommand(lp)
	handleConnectCommand := newHandlerConnectCommand(lp)

	return func(prompt string) {
		parts := strings.Split(prompt, " ")

		if len(parts) == 0 {
			logger.Error("prompt is empty")
			return
		}
		command := parts[0]
		args := parts[1:]
		switch command {
		case "serve":
			handleStartServerCommand(args)
		case "connect":
			handleConnectCommand(args)
		default:
			logger.Error("unknown command", "command", command)
		}
	}
}
