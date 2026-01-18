package mysort

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFlags парсим флаги
func ParseFlags(flags []string) (*Options, error) {
	options := &Options{}

	for i := 0; i < len(flags); i++ {
		flag := flags[i]
		if !strings.HasPrefix(flag, "-") {
			continue
		}
		flag = flag[1:] // убираем "-"

		if len(flag) >= 1 {
			// Перебираем возможные флаги, если слитно указана комбинация
			for j := 0; j < len(flag); j++ {
				switch string(flag[j]) {
				case "k":
					// Проверка есть ли после k число
					if j+1 >= len(flag) {
						return nil, fmt.Errorf("flag -k requires a value")
					}
					if len(flag) > 1 {
						column, err := strconv.Atoi(string(flag[j+1]))
						if err != nil {
							return nil, fmt.Errorf("invalid column number for -k: %s", string(flag[j+1]))
						}
						options.KeyColumn = column
						j++
					}
				case "n":
					options.Numeric = true
				case "r":
					options.Reverse = true
				case "u":
					options.Unique = true
				case "M":
					options.MonthSort = true
				case "b":
					options.IgnoreBlanks = true
				case "c":
					options.CheckOnly = true
				case "h":
					options.HumanNumeric = true

				default:
					fmt.Printf("no such flag: %s\n", flag)
				}
			}
		}
	}
	return options, nil
}
