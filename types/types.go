package types

// "STRING" represents one or multiple letters that are not keywords
const (
	Asperand     TokenVariant = "ASPERAND"
	Asterisk     TokenVariant = "ASTERISK"
	EOF          TokenVariant = "EOF"
	ForwardSlash TokenVariant = "FORWARD_SLASH"
	Hyphen       TokenVariant = "HYPHEN"
	Newline      TokenVariant = "NEWLINE"
	Number       TokenVariant = "NUMBER"
	String       TokenVariant = "STRING"
	WeightUnit   TokenVariant = "WEIGHT_UNIT"
	WhiteSpace   TokenVariant = "WHITE_SPACE"
)

type TokenVariant string

type Token struct {
	Variant TokenVariant
	Lexeme  string
	Literal interface{}
	Line    int
}

type Exercise struct {
	Name   string
	Weight Weight
	Reps   []int
}

type Weight struct {
	Value float64
	Unit  string
}
