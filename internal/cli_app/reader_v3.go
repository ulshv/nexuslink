package cli_app

// import (
// 	"bufio"
// 	"fmt"
// 	"os"

// 	"github.com/eiannone/keyboard"
// )

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
