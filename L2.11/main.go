package main

import (
	"fmt"
	"sort"
	"strings"
)

func anagram(words []string) map[string][]string {
	anagramMap := make(map[string][]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)

		// Создаем ключ - отсортированные буквы слова
		runes := []rune(lowerWord)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		key := string(runes)

		// Добавляем слово в группу анаграмм
		anagramMap[key] = append(anagramMap[key], lowerWord)
	}

	// Фильтруем группы с одной анаграммой и убираем дубликаты
	result := make(map[string][]string)
	for _, group := range anagramMap {
		if len(group) > 1 {
			// Убираем дубликаты и сортируем
			uniqueGroup := removeDuplicates(group)
			sort.Strings(uniqueGroup)
			result[uniqueGroup[0]] = uniqueGroup
		}
	}

	return result
}

func removeDuplicates(words []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}
	return result
}

func main() {
	words := []string{"пятак", "пЯтка", "тяПка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagram(words))
}
