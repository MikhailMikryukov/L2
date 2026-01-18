package main

import (
	"L2.10v2/mysort"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// Если первое слово не sort, выдаем ошибку
	command := os.Args[1]
	if command != "sort" {
		fmt.Printf("invalid command: %s", command)
		return
	}

	// Получаем данные из файла
	fileName := os.Args[len(os.Args)-1]
	fileStrings, err := getFileLines(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Обрабатываем флаги
	flags := os.Args[2 : len(os.Args)-1]
	options, err := mysort.ParseFlags(flags)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Сортируем
	sorted := mysort.ApplySortWithFlags(fileStrings, options)

	//Печать отсортированного файла
	for i := 0; i < len(sorted); i++ {
		fmt.Println(sorted[i])
	}
}

// Считываем строки из файла
func getFileLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file error: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, nil
}
