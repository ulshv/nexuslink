package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development, configure appropriately for production
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}

func getPort() string {
	// Check command line flags first
	var port string
	flag.StringVar(&port, "port", "", "Port to run the server on")
	flag.Parse()

	if port != "" {
		return ":" + port
	}

	// Check environment variable
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}

	// Default port if nothing is specified
	return ":8080"
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	port := getPort()
	fmt.Printf("WebSocket server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
