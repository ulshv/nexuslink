package cli

import "github.com/ulshv/nexuslink/pkg/log_prompt"

type CliModule struct {
	lp *log_prompt.LogPrompt
}

type cliHandlerFunc func(args []string)

func NewCliModule(lp *log_prompt.LogPrompt) *CliModule {
	return &CliModule{
		lp: lp,
	}
}
