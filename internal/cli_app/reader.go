package cli_app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	clearLine     = "\r\x1b[K" // Clear current line
	moveUpOneLine = "\x1b[1A"  // Move cursor up one line
)

type CLI struct {
	reader    *bufio.Reader
	prompt    string
	inputLine string
	mu        sync.Mutex // Add mutex for thread-safe logging
}

func NewCLI(prompt string) *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		prompt: prompt,
	}
}

func (c *CLI) ReadCommandsLoopV2(commandCh chan Command) {
	// Print initial prompt
	fmt.Print(c.prompt)

	for {
		char, _, err := c.reader.ReadRune()
		if err != nil {
			fmt.Printf("[error]: readCommandsLoop: failed to read input: %v\n", err)
			continue
		}

		switch char {
		case '\n':
			// Process the command
			fmt.Println() // Move to next line
			cleanStr := strings.TrimSpace(c.inputLine)
			params := strings.Split(cleanStr, " ")

			if len(cleanStr) > 0 {
				command := Command{
					Command: params[0],
					Args:    params[1:],
				}
				commandCh <- command
			}

			// Reset input line and show prompt
			c.inputLine = ""
			fmt.Print(c.prompt)

		case '\b', 127: // Backspace and Delete
			if len(c.inputLine) > 0 {
				c.inputLine = c.inputLine[:len(c.inputLine)-1]
				fmt.Print(clearLine + c.prompt + c.inputLine)
			}

		default:
			c.inputLine += string(char)
			fmt.Print(string(char))
		}
	}
}

// Log prints a message while preserving the current input line
func (c *CLI) Log(format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear the current line
	fmt.Print(clearLine)
	// Print the log message
	fmt.Printf(format+"\n", args...)
	// Restore prompt and current input
	fmt.Print(c.prompt + c.inputLine)
}

// func ReadCommandsLoopV1(commandCh chan Command) {
// 	for {
// 		time.Sleep(100 * time.Millisecond) // make the `>` appear after the previous log from goroutine
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("> ")
// 		inputText, err := reader.ReadString('\n')

// 		if err != nil {
// 			fmt.Printf("[error]: readCommandsLoop: failed to read input: %v\n", err)
// 			continue
// 		}

// 		cleanStr := strings.Trim(inputText, " \n")
// 		params := strings.Split(cleanStr, " ")

// 		fmt.Println("[debug]: readCommandsLoop: user input:", cleanStr)

// 		if len(params) == 0 {
// 			fmt.Println("[error]: readCommandsLoop: no command provided")
// 			continue
// 		}

// 		command := Command{
// 			Command: params[0],
// 			Args:    params[1:],
// 		}

// 		commandCh <- command
// 	}
// }
