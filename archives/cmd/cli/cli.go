package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Connect to the WebSocket server
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Handle interrupt signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a ticker that ticks every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Message counter
	counter := 0

	for {
		select {
		case <-ticker.C:
			counter++
			message := fmt.Sprintf("Hello from client! Message #%d", counter)
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}
			log.Printf("Sent: %s", message)

			// Read the response
			_, response, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Received: %s", response)

		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
