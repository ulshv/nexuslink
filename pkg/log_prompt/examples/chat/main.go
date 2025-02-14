package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ulshv/nexuslink/pkg/log_prompt"
)

var (
	usernames = []string{"alice", "bob", "olivia", "peter", "sam", "sergey", "victor"}
	messages  = []string{
		"Hello, how are you?",
		"I'm fine, thank you!",
		"What's your name?",
		"My name is anaonymous!",
		"Nice to meet you, anaonymous!",
	}
)

func randomUsername() string {
	return usernames[rand.Intn(len(usernames))]
}

func randomMessage() string {
	return messages[rand.Intn(len(messages))]
}

func main() {
	wg := &sync.WaitGroup{}
	lp := log_prompt.NewLogPrompt(context.Background(), "> ")

	go lp.Start()

	logger := lp.NewLogger("chat")
	wg.Add(2)

	go func() {
		for {
			logger.Info("Before printing", "user", "alice")
			logger.Log(fmt.Sprintf("[%s@0.0.0.0:5000]: %s", randomUsername(), randomMessage()))
			time.Sleep(time.Second)
		}
	}()

	go func() {
		currUser := "admin"

		for msg := range lp.Prompts() {
			logger.Log(fmt.Sprintf("[%s@0.0.0.0:5000]: %s", currUser, msg))
		}
	}()

	wg.Wait()
}
