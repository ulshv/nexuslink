package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// Define command line flags
	port := flag.Int("port", 8080, "Port number to listen on")
	flag.Parse()

	// Create the address string using the port
	addr := fmt.Sprintf(":%d", *port)

	// Listen for incoming connections on the specified port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", addr)

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Get client address
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("New connection from: %s\n", remoteAddr)

	// Create a buffer to store incoming data
	buffer := make([]byte, 32*1024)
	var message []byte

	for {
		// Read incoming data
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// If we have accumulated data, process it before closing
				if len(message) > 0 {
					fmt.Printf("Received from %s: %s", remoteAddr, string(message))
					_, writeErr := conn.Write(message)
					if writeErr != nil {
						log.Printf("Failed to write to connection: %v", writeErr)
					}
				}
				log.Printf("Connection closed from %s", remoteAddr)
				return
			}
			log.Printf("Error reading from %s: %v", remoteAddr, err)
			return
		}

		// Append the received data
		message = append(message, buffer[:n]...)
	}
}
