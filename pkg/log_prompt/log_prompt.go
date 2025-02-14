package log_prompt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ulshv/nexuslink/pkg/logs"
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

type logPromptLogger struct {
	*LogPrompt
	svcName string
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

func (lp *LogPrompt) NewLogger(svcName string) logs.Logger {
	return &logPromptLogger{
		LogPrompt: lp,
		svcName:   svcName,
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
	logger := lp.NewLogger("log_prompt")
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
			logger.Log("Ctrl+C received, press Ctrl+D to exit.")
			lp.currInput = ""
			lp.printPromptLine()
		case 4: // Ctrl+D
			logger.Log("Exiting the program.")
			restoreTerminalState(oldTermState) // Restore terminal state before exiting
			os.Exit(0)
		case '\n', 13: // Enter
			logger.Log(lp.prompt + lp.currInput)
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

func (l *logPromptLogger) Log(message string, args ...any) {
	if l.isLastPrompt {
		fmt.Print(CLEAR_LINE)
	}
	// get every even arg as a string to seamlesly support slog.Attr
	var params []any
	var vals []any
	for i, arg := range args {
		if i%2 == 0 {
			params = append(params, arg)
		} else {
			vals = append(vals, arg)
		}
	}
	parmsValsStringParts := []string{}
	for i, param := range params {
		parmsValsStringParts = append(parmsValsStringParts, fmt.Sprintf("%s=%+v", param, vals[i]))
	}
	logMsg := fmt.Sprintf("%s %s\n", message, strings.Join(parmsValsStringParts, ", "))
	fmt.Printf(logMsg)
	l.printPromptLine()
}

func (lpl *LogPrompt) printPromptLine() {
	fmt.Print(CLEAR_LINE)
	fmt.Print(lpl.prompt + lpl.currInput)
	lpl.isLastPrompt = true
}

// Helplers to make os.Stdin.Read() return every key stroke in the termanal:

func makeTerminalRaw() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

func restoreTerminalState(state *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), state)
}

// implements logger.Logger interface

func (l *logPromptLogger) Info(message string, args ...any) {
	l.Log(fmt.Sprintf("[INFO] [%s]: %s", l.svcName, message), args...)
}

func (l *logPromptLogger) Error(message string, args ...any) {
	l.Log(fmt.Sprintf("[ERROR] [%s]: %s", l.svcName, message), args...)
}

func (l *logPromptLogger) Warn(message string, args ...any) {
	l.Log(fmt.Sprintf("[WARN] [%s]: %s", l.svcName, message), args...)
}

func (l *logPromptLogger) Debug(message string, args ...any) {
	// return (implement debug logging later)
	l.Log(fmt.Sprintf("[DEBUG] [%s]: %s", l.svcName, message), args...)
}
