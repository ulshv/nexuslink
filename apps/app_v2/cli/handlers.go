package cli

import (
	"strings"
)

func (m *cliModule) NewCliPromptsHandler() func(prompt string) {
	logger := m.lp.NewLogger("cli/prompts_handler")

	helpCmdHandler := newHelpCmdHandler(m.lp)

	return func(prompt string) {
		if len(prompt) == 0 {
			logger.Error("prompt is empty")
		}

		parts := strings.Split(prompt, " ")
		command := parts[0]
		args := parts[1:]

		logger.Debug("received a prompt", "command", command, "args", args)

		switch command {
		case "help":
			helpCmdHandler(args...)
		default:
			logger.Error("command is not found", "command", command)
		}
	}
}
