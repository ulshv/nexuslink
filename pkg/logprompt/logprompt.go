package logprompt

import (
	"context"
	"fmt"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	CLEAR_LINE = "\r\x1b[K" // Clear current terminal line
)

type LogPrompt struct {
	prompt       string
	currInput    string
	isLastPrompt bool
	promptsCh    chan string
	ctx          context.Context
}

func NewLogPrompt(ctx context.Context, prompt string) *LogPrompt {
	return &LogPrompt{
		prompt:       prompt,
		currInput:    "",
		isLastPrompt: false,
		promptsCh:    make(chan string),
		ctx:          ctx,
	}
}

func (lp *LogPrompt) Prompts() <-chan string {
	return lp.promptsCh
}

func (lp *LogPrompt) Start() {
	// Make stdin raw mode
	oldTermState, err := makeTerminalRaw()
	if err != nil {
		fmt.Println("Failed to set raw mode:", err)
		return
	}
	// Make initial prompt line
	lp.printPromptLine()
	// Ensure we restore terminal state on exit
	defer restoreTerminalState(oldTermState)
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
		case 3: // Ctrl+C
			lp.Log("Ctrl+C received, run Ctrl+D to exit.")
			lp.currInput = ""
			lp.printPromptLine()
		case 4: // Ctrl+D
			lp.Log("Exiting the program.")
			restoreTerminalState(oldTermState) // Restore terminal state before exiting
			os.Exit(0)
		case '\n', 13: // Enter
			lp.Log(lp.prompt + lp.currInput)
			// send currInput to the channel
			lp.promptsCh <- lp.currInput
			lp.currInput = ""
			lp.printPromptLine()
		case '\b', 127: // Backspace
			if len(lp.currInput) > 0 {
				lp.currInput = lp.currInput[:len(lp.currInput)-1]
				lp.printPromptLine()
			}
		default:
			lp.currInput += string(char)
			lp.printPromptLine()
		}
	}
}

func (lp *LogPrompt) Stop() {
	lp.ctx.Done()
	// TODO: actually stop the loop started by Start()
}

func (lp *LogPrompt) Log(message string, args ...any) {
	if lp.isLastPrompt {
		fmt.Print(CLEAR_LINE)
	}
	fmt.Printf(message+"\n", args...)
	lp.printPromptLine()
}

func (lp *LogPrompt) printPromptLine() {
	fmt.Print(CLEAR_LINE)
	fmt.Print(lp.prompt + lp.currInput)
	lp.isLastPrompt = true
}

// Helpers to make os.Stdin.Read() return every key stroke in the termanal:

func makeTerminalRaw() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

func restoreTerminalState(state *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), state)
}
