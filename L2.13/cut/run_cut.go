package cut

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// RunCut - запуск реализации cut
func RunCut(config Config) error {
	//  Считываем из файла или Stdin
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

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Функция обработки пачки
	processBatch := func(lines []string) {
		defer wg.Done()
		output := make([]string, 0, len(config.Fields))

		for i := 0; i < len(lines); i++ {
			parts := strings.Split(lines[i], config.Delimiter)

			for j := 0; j < len(config.Fields); j++ {
				columnNum := config.Fields[j]
				// Если номер поля выходит за границы
				if columnNum > len(parts) || columnNum < 0 {
					continue
				}
				output = append(output, parts[columnNum])
			}
		}
		mu.Lock()
		fmt.Println(strings.Join(output, " "))
		mu.Unlock()
	}

	// Буфер на 1000 строк
	batch := make([]string, 0, 1000)

	// Считываем строки
	for scanner.Scan() {
		line := scanner.Text()

		// Игнорируем строки без разделителя (флаг -s)
		if config.Separated {
			if !strings.Contains(line, config.Delimiter) {
				continue
			}
		}

		batch = append(batch, line)

		if len(batch) == cap(batch) {
			wg.Add(1)
			go processBatch(batch)
			batch = batch[:0]
		}
	}
	// Обработка остатка
	if len(batch) > 0 {
		wg.Add(1)
		go processBatch(batch)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	return nil
}
