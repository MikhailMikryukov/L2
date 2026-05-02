package main

import (
	"L2.12/grep"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestGrepBasicFromFile(t *testing.T) {
	config := grep.Config{
		Pattern:     "hello",
		FixedString: true,
		IgnoreCase:  false,
	}

	expected := "As context for hello\nAnother hello match\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepIgnoreCaseFromFile(t *testing.T) {
	config := grep.Config{
		Pattern:     "hello",
		FixedString: true,
		IgnoreCase:  true,
	}

	expected := "Hello world this is a test\n" +
		"HELLO in uppercase\n" +
		"Hello again with context\n" +
		"As context for hello\n" +
		"Another hello match\n" +
		"Final HELLO in the file\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepWithLineNumbers(t *testing.T) {
	config := grep.Config{
		Pattern:     "Hello",
		FixedString: true,
		LineNumbers: true,
	}

	expected := "2:Hello world this is a test\n" +
		"7:Hello again with context\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepCountOnly(t *testing.T) {
	config := grep.Config{
		Pattern:     "Hello",
		FixedString: true,
		CountOnly:   true,
	}

	expected := "2\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepWithContext(t *testing.T) {
	config := grep.Config{
		Pattern:       "hello",
		FixedString:   true,
		AfterContext:  1,
		BeforeContext: 1,
	}

	expected := "This line should be captured\n" +
		"As context for hello\n" +
		"And this one too\n" +
		"Another hello match\n" +
		"With some text after\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepInvertMatch(t *testing.T) {
	config := grep.Config{
		Pattern:     "hello",
		FixedString: true,
		IgnoreCase:  true,
		InvertMatch: true,
	}

	// Все строки кроме тех, что содержат "hello" в любом регистре
	expected := "First line of the file\n" +
		"Another line here\n" +
		"Some random text\n" +
		"Goodbye world\n" +
		"This line should be captured\n" +
		"And this one too\n" +
		"With some text after\n" +
		"No match here\n" +
		"Just regular text\n" +
		"With different case\n" +
		"Last line of document\n" +
		"End of file reached\n" +
		"Final goodbye\n" +
		"The very end\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGrepWithComplexRegex(t *testing.T) {
	config := grep.Config{
		Pattern:     "\\b[hH]ello\\b", // Слово "hello" как отдельное слово (не часть другого)
		FixedString: false,
	}

	expected := "Hello world this is a test\n" +
		"Hello again with context\n" +
		"As context for hello\n" +
		"Another hello match\n"

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, output)
	}
}

func TestGrepWithRegexAndFlags(t *testing.T) {
	config := grep.Config{
		Pattern:     `[0-9]`, // Строки, содержащие цифры
		FixedString: false,
		LineNumbers: true,
		IgnoreCase:  false,
	}

	// В нашем тестовом файле нет строк с цифрами, поэтому ожидаем пустой вывод
	expected := ""

	output := runGrepWithFile(config, "test_input.txt", config.FixedString)

	if output != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, output)
	}
}

// Вспомогательная функция для запуска grep с файлом и перехватом вывода
func runGrepWithFile(config grep.Config, filename string, fixed bool) string {
	// Сохраняем оригинальные аргументы
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Устанавливаем аргументы для теста
	os.Args = []string{"grep"}

	if fixed {
		os.Args = append(os.Args, "-F")
	}

	if config.IgnoreCase {
		os.Args = append(os.Args, "-i")
	}
	if config.InvertMatch {
		os.Args = append(os.Args, "-v")
	}
	if config.CountOnly {
		os.Args = append(os.Args, "-c")
	}
	if config.LineNumbers {
		os.Args = append(os.Args, "-n")
	}
	if config.AfterContext > 0 {
		os.Args = append(os.Args, "-A", "1")
	}
	if config.BeforeContext > 0 {
		os.Args = append(os.Args, "-B", "1")
	}
	os.Args = append(os.Args, config.Pattern, filename)

	// Перехватываем stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем основную логику
	main()

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Читаем вывод
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
