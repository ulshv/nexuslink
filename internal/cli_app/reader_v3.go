package cli_app

// import (
// 	"bufio"
// 	"fmt"
// 	"os"

// 	"github.com/eiannone/keyboard"
// )

// type CLIv2 struct {
// 	reader       *bufio.Reader
// 	prompt       string
// 	currInput    string
// 	isLastPrompt bool
// }

// var cliv2 *CLIv2 = &CLIv2{
// 	reader:       bufio.NewReader(os.Stdin),
// 	prompt:       "> ",
// 	currInput:    "",
// 	isLastPrompt: false,
// }

// func processInput(input string) {
// 	fmt.Println("input:", input)
// }

// func printPromptLine() {
// 	fmt.Printf(clearLine)
// 	fmt.Printf(cliv2.prompt + cliv2.currInput)
// 	fmt.Println("cliv2.currInput:", cliv2.currInput)
// 	cliv2.isLastPrompt = true
// }

// func ReadCommandsLoopV3(commandCh chan<- Command) {
// 	for {
// 		char, key, err := keyboard.GetKey()
// 		fmt.Printf("Debug - Rune: %v, Size: %d\n", char, key)
// 		if err != nil {
// 			LogV2("[error]: readCommandsLoop: failed to read input: %v\n", err)
// 			continue
// 		} else {
// 			fmt.Println("read symbol:", char)
// 		}

// 		switch char {
// 		case '\n': // Enter key
// 			processInput(cliv2.currInput)
// 			cliv2.currInput = ""
// 		case '\b', 127: // Backspace/Delete key
// 			cliv2.currInput = cliv2.currInput[:len(cliv2.currInput)-1]
// 			printPromptLine()
// 		default: // Any other key
// 			fmt.Println("tryint to set cliv2.currInput += string(symbol)", cliv2.currInput+string(char))
// 			cliv2.currInput += string(char)
// 		}
// 	}
// }

// func LogV2(message string, args ...any) {
// 	if !cliv2.isLastPrompt {
// 		fmt.Print(clearLine)
// 	}
// 	fmt.Printf(message+"\n", args...)
// 	printPromptLine()
// }
