package cli

import "github.com/ulshv/nexuslink/pkg/log_prompt"

func newHelpCmdHandler(lp *log_prompt.LogPrompt) cmdHandler {
	logger := lp.NewLogger("cli/help_cmd_handler")

	return func(args ...string) {
		logger.Log("Welcome to the NexusLink Chat! Available commands:")
		logger.Log("help          - log this message")
		logger.Log("connect :5000 - connect to chat server on port 5000")
		logger.Log("server :5001  - start a server on port 5001")
	}
}
