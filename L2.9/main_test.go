package main

import "testing"

func TestUnwrapString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		// Базовые случаи
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "string without digits",
			input:    "abcd",
			expected: "abcd",
			wantErr:  false,
		},

		// Случаи с цифрами
		{
			name:     "single digit repetition",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			wantErr:  false,
		},
		{
			name:     "multiple digits",
			input:    "a2b3c4",
			expected: "aabbbcccc",
			wantErr:  false,
		},
		{
			name:     "digit at the beginning",
			input:    "2a",
			expected: "2a",
			wantErr:  false,
		},
		{
			name:     "digit at the end",
			input:    "a2",
			expected: "aa",
			wantErr:  false,
		},
		{
			name:     "consecutive digits",
			input:    "a2b3",
			expected: "aabbb",
			wantErr:  false,
		},

		// Случаи с экранированием
		{
			name:     "escaped digit",
			input:    "qwe\\45",
			expected: "qwe44444",
			wantErr:  false,
		},
		{
			name:     "escaped backslash",
			input:    "qwe\\\\5",
			expected: "qwe\\\\\\\\\\",
			wantErr:  false,
		},
		{
			name:     "multiple escaped characters",
			input:    "a\\2b\\3c",
			expected: "a2b3c",
			wantErr:  false,
		},

		// Ошибочные случаи
		{
			name:     "only digits",
			input:    "45",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "only digits with backslash",
			input:    "\\45",
			expected: "44444",
			wantErr:  false,
		},

		// Граничные случаи
		{
			name:     "single character",
			input:    "a",
			expected: "a",
			wantErr:  false,
		},
		{
			name:     "single digit after character",
			input:    "a1",
			expected: "a",
			wantErr:  false,
		},
		{
			name:     "max digit 9",
			input:    "a9",
			expected: "aaaaaaaaa",
			wantErr:  false,
		},
		{
			name:     "digit 2 - minimum repetition",
			input:    "a2",
			expected: "aa",
			wantErr:  false,
		},

		// Специальные случаи
		{
			name:     "mixed with spaces",
			input:    "a2 b3",
			expected: "aa bbb",
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			input:    "я2",
			expected: "яя",
			wantErr:  false,
		},
		{
			name:     "special characters",
			input:    "!2@3",
			expected: "!!@@@",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("unwrapString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("unwrapString() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestUnwrapString_EdgeCases(t *testing.T) {
	// Тест на очень длинную строку
	longInput := "a9b9c9d9e9f9"
	expected := "aaaaaaaaabbbbbbbbbcccccccccdddddddddeeeeeeeeefffffffff"
	result, err := unwrapString(longInput)
	if err != nil {
		t.Errorf("unwrapString() with long input returned error: %v", err)
	}
	if result != expected {
		t.Errorf("unwrapString() with long input = %v, expected %v", result, expected)
	}

	// Тест на строку с минимальной длиной
	result, err = unwrapString("a")
	if err != nil {
		t.Errorf("unwrapString() with single char returned error: %v", err)
	}
	if result != "a" {
		t.Errorf("unwrapString() with single char = %v, expected %v", result, "a")
	}
}
