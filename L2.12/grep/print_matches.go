package grep

import "fmt"

// PrintMatches Печать совпадений
func PrintMatches(lines []string, lineNumbers []int, matches []bool, config Config) error {
	// Отмечаем напечатанные, чтобы не повторялись строки
	printed := make([]bool, len(lines))

	// Для флага -c
	if config.CountOnly {
		matchCount := 0
		for i := 0; i < len(matches); i++ {
			if matches[i] {
				matchCount++
			}
		}
		fmt.Println(matchCount)
		return nil
	}

	// Печать совпадений
	for i := 0; i < len(lines); i++ {
		if matches[i] && !printed[i] {
			// Печать до совпадения
			start := i - config.BeforeContext
			if start < 0 {
				start = 0
			}
			for j := start; j < i; j++ {
				if !printed[j] {
					printLine(lines[j], lineNumbers[j], config.LineNumbers)
					printed[j] = true
				}
			}

			// Печать самого совпадения
			printLine(lines[i], lineNumbers[i], config.LineNumbers)
			printed[i] = true

			// Печать после совпадения
			end := i + config.AfterContext + 1
			if end > len(lines) {
				end = len(lines)
			}
			for j := i + 1; j < end; j++ {
				if !printed[j] {
					printLine(lines[j], lineNumbers[j], config.LineNumbers)
					printed[j] = true
				}
			}
		}
	}

	return nil
}

func printLine(line string, lineNum int, showNumbers bool) {
	if showNumbers {
		fmt.Printf("%d:%s\n", lineNum, line)
	} else {
		fmt.Println(line)
	}
}
