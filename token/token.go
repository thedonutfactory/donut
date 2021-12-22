package token

// Define all of monkey's tokens
const (
	Illegal = "ILLEGAL" // Token/character we don't know about
	EOF     = "EOF"     // End of file

	// Identifiers & literals
	Identifier = "IDENTIFIER" // add, foobar, x, y, ...
	Integer    = "INTEGER"
	String     = "STRING"

	// Operators
	Equal        = "="
	Plus         = "+"
	PlusPlus     = "++"
	Minus        = "-"
	MinusMinus   = "--"
	Star         = "*"
	Slash        = "/"
	Mod          = "%"
	Bang         = "!"
	EqualEqual   = "=="
	Less         = "<"
	LessEqual    = "<="
	Greater      = ">"
	GreaterEqual = ">="
	BangEqual    = "!="
	And          = "&&"
	Or           = "||"

	// Delimiters
	Comma        = ","
	Colon        = ":"
	Semicolon    = ";"
	LeftParen    = "("
	RightParen   = ")"
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
	Const    = "CONST"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"
)

// Type is a type alias for a string
type Type string

// Token is a struct representing a Monkey token - holds a Type and a Literal
type Token struct {
	Type    Type
	Literal string
	Line    int
}

var keywords = map[string]Type{
	"func":   Function,
	"let":    Let,
	"const":  Const,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// LookupIdentifier checks our keywords map for the scanned keyword. If it finds one, then
// the keywords Type is returned. If not, the user defined IDENTIFIER is returned
func LookupIdentifier(identifier string) Type {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return Identifier
}
