package cli_commands

type Command struct {
	Command string
	Args    []string
}

type CommandHandler func(args []string)

type CommandHandlers map[string]CommandHandler
