```go
package cli_app

type Command struct {
	Command string
	Args    []string
}

type CommandHandler func(cli *CLI, args []string)

type CommandHandlers map[string]CommandHandler
```
