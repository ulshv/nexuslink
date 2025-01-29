package tcp

import (
	"context"
	"fmt"

	"github.com/ulshv/nexuslink/internal/pb"
)

func HandleServerSideCommands(ctx context.Context, ch <-chan *pb.TCPCommand, conn ClientConnection) {
	for {
		select {
		case command := <-ch:
			fmt.Printf("[info]: received command on the server: %s\n", command.Command)
			switch command.Command {
			case CommandClientInit:
				SendMessage(conn, &pb.TCPCommand{
					Command: CommandServerInit,
					Payload: []byte{},
				})
			}
		case <-ctx.Done():
			return
		}
	}
}

func HandleClientSideCommands(ctx context.Context, ch <-chan *pb.TCPCommand, conn ServerConnection) {
	for {
		select {
		case command := <-ch:
			fmt.Printf("[info]: received command on the client: %s\n", command.Command)
			switch command.Command {
			case CommandServerInit:
				fmt.Println("[info]: server notified that it's initialized for the current client")
			}
		case <-ctx.Done():
			return
		}
	}
}
