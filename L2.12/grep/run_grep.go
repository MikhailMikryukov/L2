package grep

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// RunGrep Запуск утилиты grep
func RunGrep(config Config) error {

	// Установка источника откуда считываем файл или Stdin
	var input io.Reader
	if len(flag.Args()) > 1 {

		file, err := os.Open(flag.Args()[1])
		if err != nil {
			return fmt.Errorf("can't open file: %v", err)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	scanner := bufio.NewScanner(input)
	var lines []string
	var lineNumbers []int
	var matches []bool

	lineNum := 1

	// Считывание строк
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		lineNumbers = append(lineNumbers, lineNum)
		// Проверка на совпадения
		matched, err := isMatching(config.Pattern, line, config.FixedString, config.IgnoreCase)
		if err != nil {
			return err
		}

		if config.InvertMatch {
			matches = append(matches, !matched)
		} else {
			matches = append(matches, matched)
		}
		lineNum++

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
	}

	// Печать совпадений
	errPrint := PrintMatches(lines, lineNumbers, matches, config)
	if errPrint != nil {
		return errPrint
	}

	return nil
}

func isMatching(pattern string, line string, fixed bool, ignoreCase bool) (bool, error) {
	if fixed {
		if ignoreCase {
			return strings.Contains(strings.ToLower(line), strings.ToLower(pattern)), nil
		}
		return strings.Contains(line, pattern), nil

	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %v", err)
	}
	return re.MatchString(line), nil
}
