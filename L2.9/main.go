package main

import "fmt"

// Возвращает массив из N рун
func runesN(r rune, n int) []rune {
	temp := make([]rune, n)
	for i := 0; i < len(temp); i++ {
		temp[i] = r
	}
	return temp
}

// Проверяет наличие символов
func isOnlyDigits(runes []rune) bool {
	for _, r2 := range runes {
		// Если хоть один символ не цифра возвращаем false
		if r2 < '0' || r2 > '9' {
			return false
		}
	}
	return true
}

func unwrapString(s string) (string, error) {
	// Случай пустой строки
	if s == "" {
		return "", nil
	}

	runes := []rune(s)
	var result []rune

	// Случай строки только из цифр
	if isOnlyDigits(runes) {
		return "", fmt.Errorf("invalid string")
	}

	// Если первый символ цифра - просто добавляем его
	i := 0
	if runes[0] > '0' && runes[0] < '9' {
		result = append(result, runes[i])
		i++
	}

	for ; i < len(runes); i++ {
		switch runes[i] {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n := int(runes[i]) - 49
			result = append(result, runesN(runes[i-1], n)...)
		case '\\':
			i++
			result = append(result, runes[i])
		default:
			result = append(result, runes[i])
		}
	}
	return string(result), nil
}

func main() {
	fmt.Println(unwrapString("a4bc2d5e"))
	fmt.Println(unwrapString("abcd"))
	fmt.Println(unwrapString("45"))
	fmt.Println(unwrapString(""))
	fmt.Println(unwrapString("qwe\\45"))
	fmt.Println(unwrapString("qwe\\4\\5"))
	fmt.Println(unwrapString("a\\2b\\3c"))
}
