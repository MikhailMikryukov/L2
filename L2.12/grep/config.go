package grep

// Config флаги и паттерн
type Config struct {
	Pattern       string
	AfterContext  int
	BeforeContext int
	CountOnly     bool
	IgnoreCase    bool
	InvertMatch   bool
	FixedString   bool
	LineNumbers   bool
}
