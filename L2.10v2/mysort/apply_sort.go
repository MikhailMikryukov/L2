package mysort

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

var dateFormat = "02-Jan-2006"

// ApplySortWithFlags сортируем по указанным флагам
func ApplySortWithFlags(lines []string, options *Options) []string {
	// Проверка сортировки ДО любой обработки
	if options.CheckOnly {
		if isSorted(lines) {
			return nil // файл уже отсортирован
		}
	}

	// Уникальность применяется перед сортировкой
	if options.Unique {
		lines = removeDuplicates(lines)
	}

	// Игнорируем хвостовые пробелы
	if options.IgnoreBlanks {
		lines = ignoreTrailingBlanks(lines)
	}

	// Сортировка по столбцу
	if options.KeyColumn > 0 {
		lines = sortWithFlagK(options.KeyColumn, lines)
	}

	// Сортировка по числовому значению
	if options.Numeric {
		lines = numericSort(lines)
	}

	// Сортировка с учетом суффиксов
	if options.HumanNumeric {
		lines = humanReadableSort(lines)
	}

	// Сортировка по названию месяца
	if options.MonthSort {
		lines = monthSort(lines)
	}

	// Сортировка в обратном порядке
	if options.Reverse {
		lines = reverse(lines)
	}

	return lines
}

// Проверка отсортирован ли массив
func isSorted(lines []string) bool {
	for i := 0; i < len(lines); i++ {
		if lines[i+1] < lines[i] {
			return false
		}
	}
	return true
}

// Сортировка в обратном порядке
func reverse(arrToSort []string) []string {
	for i := 0; i < len(arrToSort)/2; i++ {
		arrToSort[i], arrToSort[len(arrToSort)-1-i] = arrToSort[len(arrToSort)-1-i], arrToSort[i]
	}
	return arrToSort
}

// Удаление повторяющихся строк
func removeDuplicates(arrToSort []string) []string {
	set := make(map[string]struct{})
	var unique []string

	// Проверяем есть ли строка в мапе, если нет добавляем
	for _, str := range arrToSort {
		if _, ok := set[str]; !ok {
			set[str] = struct{}{}
			unique = append(unique, str)
		}
	}

	return unique
}

// Сортировка по числовому значению
func numericSort(arrToSort []string) []string {
	_, err := strconv.Atoi(arrToSort[0])
	if err != nil {
		sort.Strings(arrToSort)
	} else {
		strAsInt := make([]int, len(arrToSort))
		for i := 0; i < len(arrToSort); i++ {
			strAsInt[i], err = strconv.Atoi(arrToSort[i])
			if err != nil {
				log.Fatal(err)
			}
		}
		sort.Ints(strAsInt)
		for i := 0; i < len(arrToSort); i++ {
			arrToSort[i] = strconv.Itoa(strAsInt[i])
		}
	}
	return arrToSort
}

// Сортировка по числовому значению
func sortWithFlagK(column int, arrToSort []string) []string {
	// Если номер колонки больше чем количество колонок, то не сортируем
	if column > len(arrToSort) || column < 0 {
		return arrToSort
	}
	// В мапе ключи - данные из колонки, значения - вся строка
	m := make(map[string]string)
	// В этом срезе будем сортировать данные из колонки
	columnsToSort := make([]string, len(arrToSort))
	for i := 0; i < len(arrToSort); i++ {
		//разделяем строку
		fields := strings.Fields(arrToSort[i])
		// кладем в мапу пару ключ - значение
		m[fields[column-1]] = arrToSort[i]
		// заполняем данные из колонок
		columnsToSort[i] = fields[column-1]
	}
	// Сортируем колонки
	sortedColumns := numericSort(columnsToSort)
	// Выставляем в нужном порядке строки
	for i := 0; i < len(sortedColumns); i++ {
		arrToSort[i] = m[sortedColumns[i]]
	}
	return arrToSort
}

// Сортировка по числовому значению (человеко-читаемые размеры)
func humanReadableSort(arrToSort []string) []string {
	// Мап ключ - переведенная в инт строка из arrToSort, значение - сама эта строка
	m := make(map[int]string)
	// Срез для численной сортировки
	temp := make([]int, len(arrToSort))

	for i := 0; i < len(arrToSort); i++ {
		stringNumber := arrToSort[i]
		// Перевод чисел с суффиксами
		if strings.HasSuffix(stringNumber, "K") || strings.HasSuffix(stringNumber, "M") ||
			strings.HasSuffix(stringNumber, "G") {
			num, err := strconv.ParseFloat(stringNumber[:len(stringNumber)-1], 32)
			if err != nil {
				fmt.Println(err)
			}
			suffix := stringNumber[len(stringNumber)-1:]
			if suffix == "K" {
				temp[i] = int(num * 1024)
				m[temp[i]] = stringNumber
			}
			if suffix == "M" {
				temp[i] = int(num * 1024 * 1024)
				m[temp[i]] = stringNumber
			}

			if suffix == "G" {
				temp[i] = int(num * 1024 * 1024 * 1024)
				m[temp[i]] = stringNumber
			}

		} else {
			// Перевод чисел без суффиксов
			var err error
			temp[i], err = strconv.Atoi(stringNumber)
			if err != nil {
				fmt.Println(err)
			}
			m[temp[i]] = stringNumber
		}
	}
	// Сортировка
	sort.Ints(temp)

	sorted := make([]string, len(temp))
	for i := 0; i < len(sorted); i++ {
		sorted[i] = m[temp[i]]
	}
	return sorted
}

// Сортировка по названию месяца
func monthSort(arrToSort []string) []string {
	parsed := make([]time.Time, len(arrToSort))
	// Мап ключ - время, значение соответсвующая строка
	m := make(map[time.Time]string)

	for i := 0; i < len(arrToSort); i++ {
		var err error
		parsed[i], err = time.Parse(dateFormat, arrToSort[i])
		// При ошибке парсинга устанавливаем большую дату, чтобы закинуть в конец списка при сортировке
		if err != nil {
			fmt.Println(err)
			parsed[i] = time.Date(9990+i, 12, 31, 0, 0, 0, 0, time.UTC)
		}
		m[parsed[i]] = arrToSort[i]

	}
	// Сортируем время
	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].Before(parsed[j])
	})
	// Выставляем в отсортированном порядке
	for i := 0; i < len(parsed); i++ {
		arrToSort[i] = m[parsed[i]]
	}

	return arrToSort
}

// Игнорировать хвостовые пробелы
func ignoreTrailingBlanks(arrToSort []string) []string {
	result := make([]string, len(arrToSort))

	for i := 0; i < len(arrToSort); i++ {
		strWithoutTrailBlanks := arrToSort[i]
		// Перебираем символы строки с конца на наличие пробелов
		for j := len(arrToSort[i]) - 1; j >= 0; j-- {
			if string(arrToSort[i][j]) == " " {
				strWithoutTrailBlanks = arrToSort[i][0:j]
			} else {
				// Если пробелов не осталось, то идем дальше
				break
			}
		}
		result[i] = strWithoutTrailBlanks
	}

	return result
}
