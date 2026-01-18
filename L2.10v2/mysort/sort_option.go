package mysort

// Options Записываем в структуру флаги, по которым нужно сортировать
type Options struct {
	KeyColumn    int  // -kN
	Numeric      bool // -n
	Reverse      bool // -r
	Unique       bool // -u
	MonthSort    bool // -M
	IgnoreBlanks bool // -b
	CheckOnly    bool // -c
	HumanNumeric bool // -h
}
