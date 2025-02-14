TCP server:
```go
// start listener
listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", host, port))
// ...
for {
  // accept connections
  conn, err := listener.Accept()
  // ...
}
```

TCP client:
```go
// connect to a server
conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
// write a message to the server
_, err := conn.Write([]byte(msg))
```

Starting the log_prompt
```go
wg := sync.WaitGroup{}
	lp := log_prompt.NewLogPrompt(context.Background(), "> ")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for prompt := range lp.Prompts() {
			handleCliCommands(lp, prompt)
		}
	}()

	lp.Start()
	wg.Wait()
```
