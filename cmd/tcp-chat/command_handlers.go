package main

import (
	"fmt"

	"github.com/ulshv/nexuslink/pkg/tcp"
)

func helloHandler(args []string) {
	fmt.Println(`Available commands:
- start <port>: Start a server
- connect <port> [<host> (optional, default 'localhost')]: Connect to a server
- [TODO] nickname <server_uuid> <nickname>: Reserve a nickname on the server
- [TODO] create_room <server_uuid> <room_name>: Create a room in the server
- [TODO] list: List all connected servers
- [TODO] disconnect <server_uuid>: Disconnect from the server
- [TODO] info <server_uuid>: Show info about a server
- [TODO] list_rooms <server_uuid>: List all rooms in the server
- [TODO] join_room <server_uuid> <room_name>: Join a room in the server
- [TODO] leave_room <server_uuid> <room_name>: Leave a room in the server
- [TODO] send_message <server_uuid> <room_name> <message...>: Send a message to the server
- [TODO] send <message>: Send messge to the latest server/room
- [TODO] exit: Exit the program`)
}

func startHandler(args []string) {
	if len(args) != 1 {
		fmt.Println("[error]: start: invalid number of arguments (need <port>)")
		return
	}
	fmt.Println("[info]: starting TCP server")
	server := tcp.NewServer(&tcp.NewServerConfig{
		Host: "0.0.0.0",
		Port: args[0],
	})
	go server.RunServer()
}

func connectHandler(args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("[error]: connect: invalid number of arguments (need <port> [<host> optional])")
		return
	}
	fmt.Printf("[info]: connecting to the server on port %s...\n", args[0])
	port := args[0]
	host := "localhost"
	if len(args) == 2 {
		host = args[1]
	}
	client, err := tcp.NewClient(tcp.NewClientConfig{
		Address: host,
		Port:    port,
	})
	if err != nil {
		fmt.Println("[error]: connect: failed to connect to the server: ", err)
		return
	}
	go client.SendMessage("Hello, world!")
}
