```go
package tcp

import (
	"context"
	"fmt"

	"github.com/ulshv/nexuslink/internal/pb"
)

func HandleServerSideCommands(ctx context.Context, ch <-chan *pb.TCPMessage, conn ClientConnection) {
	for {
		select {
		case command := <-ch:
			fmt.Printf("[info]: received command on the server: %s\n", command.Type)
			switch command.Type {
			case CommandClientInit:
				PublishMessage(conn, &pb.TCPMessage{
					Type:    CommandServerInit,
					Payload: []byte{},
				})
			}
		case <-ctx.Done():
			return
		}
	}
}

func HandleClientSideCommands(ctx context.Context, ch <-chan *pb.TCPMessage, conn ServerConnection) {
	for {
		select {
		case command := <-ch:
			fmt.Printf("[info]: received command on the client: %s\n", command.Type)
			switch command.Type {
			case CommandServerInit:
				fmt.Println("[info]: server notified that it's initialized for the current client")
			}
		case <-ctx.Done():
			return
		}
	}
}
```
