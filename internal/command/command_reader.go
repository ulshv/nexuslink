package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func ReadCommandsLoop(commandCh chan Command) {
	for {
		time.Sleep(100 * time.Millisecond) // make the `>` appear after the previous log from goroutine
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		inputText, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("[error]: readCommandsLoop: failed to read input: %v\n", err)
			continue
		}

		cleanStr := strings.Trim(inputText, " \n")
		params := strings.Split(cleanStr, " ")

		fmt.Println("[debug]: readCommandsLoop: user input:", cleanStr)

		if len(params) == 0 {
			fmt.Println("[error]: readCommandsLoop: no command provided")
			continue
		}

		command := Command{
			Command: params[0],
			Args:    params[1:],
		}

		commandCh <- command
	}
}
