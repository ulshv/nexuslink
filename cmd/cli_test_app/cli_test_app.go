package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/ulshv/nexuslink/internal/cli_app"
)

func main() {
	// ch := make(chan cli_app.Command)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	go func() {
		counter := 0
		for {
			fmt.Println("counter:", counter)
			// cli_app.LogV2("[debug] test cli log, counter: %v", counter)
			counter++
			time.Sleep(1 * time.Second)
		}
	}()

	go cli_app.ReadCommandsLoopV4()

	wg.Wait()
}
