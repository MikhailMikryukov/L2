package minishell

import (
	"bufio"
	"os"
)

func Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()
		execCommand(command)
	}
}
