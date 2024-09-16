package token

type Type string

const (
	Whitespace Type = "whitespace"
	Separator  Type = "separator"
	Bool       Type = "bool"
	String     Type = "string"
	Number     Type = "number"
	Null       Type = "null"
	EOF        Type = "eof"
	UNKNOWN    Type = "unknown"
)

type Token struct {
	Type  Type
	Value string
}
