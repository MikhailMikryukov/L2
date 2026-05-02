package minishell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var currentDir = "C:"

func execCommand(command string) {
	fields := strings.Fields(command)

	switch fields[0] {
	case "echo":
		strings.Join(fields[1:], " ")
	case "changeDirectory":
		err := changeDirectory(fields)
		if err != nil {
			fmt.Println(err)
		}
	case "pwd":
		fmt.Printf("Current folder: %s", currentDir)
	case "kill":
		err := killProcessByPID(fields[1])
		if err != nil {
			fmt.Println(err)
		}
	case "ps":
		err := ps()
		if err != nil {
			fmt.Println(err)
		}

	}
}

func changeDirectory(args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("changeDirectory: too many arguments")
	}
	dir := args[1]
	temp := path.Join(currentDir, dir)
	info, err := os.Stat(temp)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory")
	}

	currentDir = temp

	return nil
}

func ps() error {
	cmd := exec.Command("tasklist")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	return nil
}

func killProcessByPID(pid string) error {
	cmd := exec.Command("taskkill", "/PID", pid, "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", output)
	}
	fmt.Println(string(output))
	return nil
}
