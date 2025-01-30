package cli_app

import (
	"fmt"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

type CLIv2 struct {
	prompt       string
	currInput    string
	isLastPrompt bool
}

var cliv2 *CLIv2 = &CLIv2{
	prompt:       "$> ",
	currInput:    "",
	isLastPrompt: false,
}

func ReadCommandsLoopV4() {
	// Make stdin raw mode
	oldState, err := makeTerminalRaw()
	if err != nil {
		fmt.Println("Failed to set raw mode:", err)
		return
	}
	// Ensure we restore terminal state on exit
	defer restoreTerminalState(oldState)
	// Buffer for UTF-8/32 bit characters
	buf := make([]byte, 4)
	for {
		// Read key strokes on terminal
		n, err := os.Stdin.Read(buf)
		char, _ := utf8.DecodeRune(buf[:n])
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		// Handle key strokes
		switch char {
		case 3, 4: // Ctrl+C, Ctrl+D
			restoreTerminalState(oldState) // Restore terminal state before exiting
			LogV2("Exiting the program.")
			os.Exit(0)
		case '\n', 13: // Enter
			LogV2(cliv2.prompt + cliv2.currInput)
			// TODO: send currInput to `chan<- string`
			// processInput(cliv2.currInput)
			cliv2.currInput = ""
			printPromptLine()
		case '\b', 127: // Backspace
			if len(cliv2.currInput) > 0 {
				cliv2.currInput = cliv2.currInput[:len(cliv2.currInput)-1]
				printPromptLine()
			}
		default:
			cliv2.currInput += string(char)
			printPromptLine()
		}
	}
}

func LogV2(message string, args ...any) {
	if cliv2.isLastPrompt {
		fmt.Print(clearLine)
	}
	fmt.Printf(message+"\n", args...)
	printPromptLine()
}

func printPromptLine() {
	fmt.Print(clearLine)
	fmt.Printf(cliv2.prompt + cliv2.currInput)
	cliv2.isLastPrompt = true
}

// Helpers to make os.Stdin.Read() return every key stroke in the termanal:

func makeTerminalRaw() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

func restoreTerminalState(state *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), state)
}
