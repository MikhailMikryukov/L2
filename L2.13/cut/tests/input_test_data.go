package tests

import (
	"L2.13/cut"
	"strings"
)

// TestInputData представляет тестовые входные данные
type TestInputData struct {
	Name     string
	Input    string
	Config   cut.Config
	Expected string
}

// GetTestCases возвращает все тестовые случаи для входных данных
func GetTestCases() []TestInputData {
	return []TestInputData{
		// Базовые тесты
		{
			Name:  "basic comma separated",
			Input: "apple,banana,cherry,date\none,two,three,four",
			Config: cut.Config{
				Fields:    []int{1, 3},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "banana date\ntwo four",
		},
		{
			Name:  "single field",
			Input: "name:age:city:country\nJohn:25:London:UK",
			Config: cut.Config{
				Fields:    []int{2},
				Delimiter: ":",
				Separated: false,
			},
			Expected: "city\nLondon",
		},
		{
			Name:  "tab delimiter",
			Input: "col1\tcol2\tcol3\tcol4\nval1\tval2\tval3\tval4",
			Config: cut.Config{
				Fields:    []int{0, 2},
				Delimiter: "\t",
				Separated: false,
			},
			Expected: "col1 col3\nval1 val3",
		},

		// Тесты с флагом -s (separated)
		{
			Name:  "separated flag with delimiter lines",
			Input: "a,b,c\nno-delimiter\nx,y,z\nanother-no-delimiter",
			Config: cut.Config{
				Fields:    []int{1},
				Delimiter: ",",
				Separated: true,
			},
			Expected: "b\ny",
		},
		{
			Name:  "separated flag without delimiter lines",
			Input: "a,b,c\nx,y,z\np,q,r",
			Config: cut.Config{
				Fields:    []int{0},
				Delimiter: ",",
				Separated: true,
			},
			Expected: "a\nx\np",
		},

		// Тесты с несуществующими полями
		{
			Name:  "field out of range",
			Input: "one,two,three\n1,2,3",
			Config: cut.Config{
				Fields:    []int{5},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "\n",
		},
		{
			Name:  "mixed existing and non-existing fields",
			Input: "a,b,c,d,e\n1,2,3,4,5",
			Config: cut.Config{
				Fields:    []int{0, 10, 2},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "a c\n1 3",
		},

		// Тесты с отрицательными полями
		{
			Name:  "negative field index",
			Input: "one,two,three\n1,2,3",
			Config: cut.Config{
				Fields:    []int{-1},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "\n",
		},

		// Тесты с пустыми полями
		{
			Name:  "empty fields in middle",
			Input: "a,,c\n1,,3",
			Config: cut.Config{
				Fields:    []int{0, 1, 2},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "a  c\n1  3",
		},

		// Тесты с разными разделителями
		{
			Name:  "semicolon delimiter",
			Input: "first;second;third\nA;B;C",
			Config: cut.Config{
				Fields:    []int{0, 2},
				Delimiter: ";",
				Separated: false,
			},
			Expected: "first third\nA C",
		},
		{
			Name:  "space delimiter",
			Input: "apple banana cherry\nred green blue",
			Config: cut.Config{
				Fields:    []int{1},
				Delimiter: " ",
				Separated: false,
			},
			Expected: "banana\ngreen",
		},

		// Тесты с одной строкой
		{
			Name:  "single line input",
			Input: "alpha,beta,gamma,delta",
			Config: cut.Config{
				Fields:    []int{1, 3},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "beta delta",
		},

		// Тесты с пустым вводом
		{
			Name:  "empty input",
			Input: "",
			Config: cut.Config{
				Fields:    []int{1, 2},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "",
		},

		// Тесты только с разделителями
		{
			Name:  "only delimiters",
			Input: ",,,\n,,,",
			Config: cut.Config{
				Fields:    []int{1},
				Delimiter: ",",
				Separated: false,
			},
			Expected: " \n ",
		},

		// Тесты с полями в разном порядке
		{
			Name:  "fields in reverse order",
			Input: "one,two,three,four\n1,2,3,4",
			Config: cut.Config{
				Fields:    []int{3, 1, 0},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "four two one\n4 2 1",
		},

		// Тесты с дублирующимися полями
		{
			Name:  "duplicate fields",
			Input: "a,b,c,d\n1,2,3,4",
			Config: cut.Config{
				Fields:    []int{1, 1, 1},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "b b b\n2 2 2",
		},
	}
}

// GetEdgeCases возвращает граничные случаи
func GetEdgeCases() []TestInputData {
	return []TestInputData{
		{
			Name:  "very long line",
			Input: "a," + strings.Repeat("x,", 1000) + "z\nshort,line",
			Config: cut.Config{
				Fields:    []int{0, 1001},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "a z\nshort\n",
		},
		{
			Name:  "empty lines with content",
			Input: "\n\na,b,c\n\nx,y,z\n\n",
			Config: cut.Config{
				Fields:    []int{1},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "\n\nb\n\ny\n\n",
		},
		{
			Name:  "only newlines",
			Input: "\n\n\n",
			Config: cut.Config{
				Fields:    []int{1},
				Delimiter: ",",
				Separated: false,
			},
			Expected: "\n\n\n",
		},
	}
}
