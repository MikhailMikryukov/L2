package main

import (
	"L2.12/grep"
	"fmt"
	"os"
)

func main() {
	config := grep.ParseFlags()

	err := grep.RunGrep(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v", err)
	}
}
