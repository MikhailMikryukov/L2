package tests

import (
	"L2.13/cut"
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestRunCutWithTestData(t *testing.T) {
	testCases := GetTestCases()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Создаем буфер для ввода
			input := strings.NewReader(tc.Input)

			// Создаем буфер для вывода
			var output bytes.Buffer

			// Сохраняем оригинальный stdout и заменяем на наш буфер
			oldStdout := stdout
			stdout = &output
			defer func() { stdout = oldStdout }()

			// Запускаем cut с тестовыми данными
			err := runCutWithConfig(input, tc.Config)
			if err != nil {
				t.Errorf("RunCut() error = %v", err)
				return
			}

			// Получаем результат
			actual := strings.TrimSpace(output.String())
			expected := strings.TrimSpace(tc.Expected)

			if actual != expected {
				t.Errorf("RunCut() = %q, want %q", actual, expected)
			}
		})
	}
}

func TestRunCutEdgeCases(t *testing.T) {
	edgeCases := GetEdgeCases()

	for _, tc := range edgeCases {
		t.Run(tc.Name, func(t *testing.T) {
			input := strings.NewReader(tc.Input)
			var output bytes.Buffer

			oldStdout := stdout
			stdout = &output
			defer func() { stdout = oldStdout }()

			err := runCutWithConfig(input, tc.Config)
			if err != nil {
				t.Errorf("RunCut() error = %v", err)
				return
			}

			actual := output.String()
			expected := tc.Expected

			if actual != expected {
				t.Errorf("RunCut() = %q, want %q", actual, expected)
			}
		})
	}
}

// runCutWithConfig - вспомогательная функция для тестирования
func runCutWithConfig(input *strings.Reader, config cut.Config) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		// Игнорируем строки без разделителя если установлен флаг -s
		if config.Separated && !strings.Contains(line, config.Delimiter) {
			continue
		}

		parts := strings.Split(line, config.Delimiter)
		var outputFields []string

		for _, field := range config.Fields {
			if field >= 0 && field < len(parts) {
				outputFields = append(outputFields, parts[field])
			}
		}

		if len(outputFields) > 0 {
			outputLine := strings.Join(outputFields, " ")
			stdout.Write([]byte(outputLine + "\n"))
		} else {
			stdout.Write([]byte("\n"))
		}
	}

	return scanner.Err()
}

// Глобальная переменная для stdout (для тестирования)
var stdout = &bytes.Buffer{}
