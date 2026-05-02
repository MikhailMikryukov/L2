package cut

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// ParseFlags - обработка флагов
func ParseFlags() (Config, error) {
	config := Config{}
	var fieldsStr string

	flag.StringVar(&fieldsStr, "f", "", "Print the specified fields")
	flag.StringVar(&config.Delimiter, "d", "\\t", "Set the delimiter")
	flag.BoolVar(&config.Separated, "s", false, "Ignore lines without delimiter")

	flag.Parse()

	fields, err := parseFieldsStr(fieldsStr)
	if err != nil {
		return config, err
	}
	config.Fields = fields

	return config, nil
}

func parseFieldsStr(fieldsStr string) ([]int, error) {
	// Случай если не указан флаг или поля
	if len(fieldsStr) < 1 {
		return nil, nil
	}

	fields := strings.Split(fieldsStr, ",")
	result := make([]int, 0, len(fields))

	for i := 0; i < len(fields); i++ {
		numStr := fields[i]

		// Обработка если указан диапазон полей
		if strings.Contains(numStr, "-") {
			firstNum, err1 := strconv.Atoi(string(numStr[0]))
			secondNum, err2 := strconv.Atoi(string(numStr[2]))
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range: %v", err1)
			}

			for j := firstNum; j <= secondNum; j++ {
				result = append(result, j)
			}
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("not a number: %v", err)
		}

		result = append(result, num)
	}

	return result, nil
}
