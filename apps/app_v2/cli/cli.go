package cli

import "github.com/ulshv/nexuslink/pkg/log_prompt"

type cliModule struct {
	lp *log_prompt.LogPrompt
}

type cmdHandler func(args ...string)

func NewCliModule(lp *log_prompt.LogPrompt) *cliModule {
	return &cliModule{lp}
}
