package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/ulshv/nexuslink/pkg/logprompt"
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
	lp := logprompt.NewLogPrompt(context.Background(), "> ")

	wg.Add(2)

	go func() {
		for {
			lp.Log("[%s@0.0.0.0:5000]: %s", randomUsername(), randomMessage())
			time.Sleep(time.Second)
		}
	}()

	go func() {
		currUser := "admin"

		for msg := range lp.Prompts() {
			lp.Log("[%s@0.0.0.0:5000]: %s", currUser, msg)
		}
	}()

	lp.Start()
}
