package grep

import (
	"flag"
)

// ParseFlags Парсинг флагов
func ParseFlags() Config {
	config := Config{}
	var cFlagCount int

	flag.IntVar(&config.AfterContext, "A", 0, "Print N strings after match")
	flag.IntVar(&config.BeforeContext, "B", 0, "Print N strings before match")
	flag.IntVar(&cFlagCount, "C", 0, "Print N strings before and after match")
	flag.BoolVar(&config.CountOnly, "c", false, "Count matching lines only")
	flag.BoolVar(&config.IgnoreCase, "i", false, "Ignore case")
	flag.BoolVar(&config.InvertMatch, "v", false, "Invert match")
	flag.BoolVar(&config.FixedString, "F", false, "Fixed string match")
	flag.BoolVar(&config.LineNumbers, "n", false, "Print line numbers")

	flag.Parse()

	if cFlagCount > 0 {
		config.BeforeContext = cFlagCount
		config.AfterContext = cFlagCount
	}

	args := flag.Args()
	if len(args) > 1 {
		config.Pattern = args[0]
	}

	return config
}
