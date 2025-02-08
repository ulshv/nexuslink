package main

import (
	"strconv"
	"strings"

	"github.com/ulshv/nexuslink/pkg/tcp/tcp_client"
)

func handleConnect(client *ChatClient, params []string) {
	if len(params) == 0 {
		client.lp.Log("connect: missing params. Usage: 'connect <host:port>'")
		return
	}
	hostPort := strings.Split(params[0], ":")
	if len(hostPort) != 2 {
		client.lp.Log("connect: invalid host:port. Usage: 'connect <host:port>', 'connect %s' provided", params[0])
		return
	}
	host := hostPort[0]
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		client.lp.Log("connect: invalid port. Should be a number, '%s' provided", hostPort[1])
		return
	}
	client.lp.Log("Connecting to %s:%v...", host, port)
	serverConn, err := tcp_client.NewServerConnection(host, port)
	if err != nil {
		client.lp.Log("connect: failed to connect to the server: %v", err)
		return
	}
	client.lp.Log("Connected to the server")
	client.setServerConn(serverConn)
	client.sendMessage([]byte("what's up!"))
}
