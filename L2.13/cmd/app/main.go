package main

import (
	"L2.13/cut"
	"fmt"
)

func main() {
	config, err := cut.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cut.RunCut(config)
	if err != nil {
		fmt.Println(err)
		return
	}
}
