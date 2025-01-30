package cli_app

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ReadCommandsLoopV4() {
	readStdinOneByOne()
}

func readStdinOneByOne() {
	// Make stdin raw mode
	oldState, err := makeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to set raw mode:", err)
		return
	}
	defer restore(oldState) // Ensure we restore terminal state on exit

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		// Handle special keys
		switch buf[0] {
		case 3, 4: // Ctrl+C, Ctrl+D
			fmt.Println("^C")
			restore(oldState) // Restore terminal state before exiting
			fmt.Println("Exiting the program.")
			os.Exit(0)
		case 13: // Enter
			processInput(cliv2.currInput)
			cliv2.currInput = ""
			printPromptLine()
		case 127: // Backspace
			if len(cliv2.currInput) > 0 {
				cliv2.currInput = cliv2.currInput[:len(cliv2.currInput)-1]
				printPromptLine()
			}
		default:
			if buf[0] >= 32 { // Printable characters
				cliv2.currInput += string(buf[0])
				printPromptLine()
			}
		}
	}
}

// Add these new functions at the end of the file
func makeRaw(fd int) (*term.State, error) {
	return term.MakeRaw(fd)
}

func restore(state *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), state)
}
