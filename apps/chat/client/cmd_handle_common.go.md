```go
package main

func handleUnknownCommand(client *ChatClient, command string) {
	client.lp.Log("Unknown command: %s. Type 'help' to see available commands.", command)
}

func handleHelp(client *ChatClient) {
	client.lp.Log("Available commands:")
	client.lp.Log("- connect <host:port>: Connect to a server")
	client.lp.Log("  examples: 'connect :5000' (localhost), 'connect chat.nexuslink.dev:5000' (remote server)")
}
```
