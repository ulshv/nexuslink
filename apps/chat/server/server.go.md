```go
package main

import (
	"flag"
	"log/slog"

	tcp_server "github.com/ulshv/nexuslink/internal/tcp/tcp_server"
)

func main() {
	port := flag.Int("p", 5000, "port to listen on")
	flag.Parse()
	slog.Info("starting server...", "port", *port)
	server := tcp_server.NewServer("0.0.0.0", *port)
	server.ListenAndHandle()
}
```
