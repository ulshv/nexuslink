package logger

import "fmt"

func СlearCurrentLine() {
	fmt.Print("\n\033[1A\033[K")
}
