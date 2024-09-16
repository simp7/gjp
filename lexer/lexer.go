package lexer

import (
	"io"
	"unicode"

	"github.com/simp7/gjp/token"
)

type Lexer struct {
	text    []rune
	current int
	ch      rune
}

func NewLexer(input string) *Lexer {
	l := new(Lexer)
	l.text = []rune(input)
	return l
}

func (l *Lexer) readCh() error {
	if len(l.text) == l.current {
		return io.EOF
	}
	l.ch = l.text[l.current]
	l.current++
	return nil
}

func (l *Lexer) peakCh() rune {
	return l.text[l.current+1]
}

func (l *Lexer) parseString() token.Token {
	l.readCh()
	value := make([]rune, 0)
	for !(l.ch == '"') {
		if l.ch == '\\' {
			value = append(value, '\\')
			l.readCh()
			special := l.ch
			if special == '"' || special == '\\' || special == '/' || special == 'b' || special == 'f' || special == 'n' || special == 'r' || special == 't' {
				value = append(value, special)
			} else if special == 'u' {
				for i := 0; i < 4; i++ {
					l.readCh()
					if !unicode.In(l.ch, unicode.Hex_Digit) {
						return token.Token{Type: token.UNKNOWN, Value: string(value)}
					}
					value = append(value, l.ch)
				}
			} else {
				return token.Token{Type: token.UNKNOWN, Value: string(value)}
			}
			l.readCh()
			continue
		}
		value = append(value, l.ch)
		if err := l.readCh(); err != nil {
			return token.Token{Type: token.UNKNOWN, Value: string(value)}
		}
	}
	l.readCh()
	return token.Token{Type: token.String, Value: string(value)}
}

func (l *Lexer) parseKeyword(tokenType token.Type, keyword string) token.Token {
	value := make([]rune, 0)
	for _, v := range keyword {
		value = append(value, l.ch)
		if l.ch != v || l.readCh() != nil {
			return token.Token{Type: token.UNKNOWN, Value: string(value)}
		}
	}
	return token.Token{Type: tokenType, Value: string(value)}
}

func (l *Lexer) parseNumber() token.Token {
	value := make([]rune, 0)
	if l.ch == '-' {
		value = append(value, l.ch)
		l.readCh()
	}
	if !unicode.IsDigit(l.ch) || (l.ch == '0' && unicode.IsDigit(l.peakCh())) {
		return token.Token{Type: token.UNKNOWN, Value: string(value)}
	}
	for unicode.IsDigit(l.ch) {
		value = append(value, l.ch)
		l.readCh()
	}
	if l.ch == '.' {
		value = append(value, l.ch)
		if !unicode.IsDigit(l.peakCh()) {
			return token.Token{Type: token.UNKNOWN, Value: string(value)}
		}
		l.readCh()
		for unicode.IsDigit(l.ch) {
			value = append(value, l.ch)
			l.readCh()
		}
	}
	if l.ch == 'e' || l.ch == 'E' {
		value = append(value, l.ch)
		l.readCh()
		if l.ch == '+' || l.ch == '-' {
			value = append(value, l.ch)
			l.readCh()
		}

		if !unicode.IsDigit(l.ch) {
			return token.Token{Type: token.UNKNOWN, Value: string(value)}
		}
		for unicode.IsDigit(l.ch) {
			value = append(value, l.ch)
			l.readCh()
		}
	}

	return token.Token{Type: token.Number, Value: string(value)}
}

func (l *Lexer) parseWhitespace() token.Token {
	value := make([]rune, 0)
	for unicode.IsSpace(l.peakCh()) {
		value = append(value, l.ch)
		l.readCh()
	}
	value = append(value, l.ch)
	return token.Token{Type: token.Whitespace, Value: string(value)}
}

func (l *Lexer) NextToken() token.Token {
	if len(l.text) == l.current {
		return token.Token{Type: token.EOF, Value: ""}
	}

	if err := l.readCh(); err != nil {
		return token.Token{Type: token.EOF}
	}
	switch l.ch {
	case '{':
		return token.Token{Type: token.Separator, Value: "{"}
	case '}':
		return token.Token{Type: token.Separator, Value: "}"}
	case '[':
		return token.Token{Type: token.Separator, Value: "["}
	case ']':
		return token.Token{Type: token.Separator, Value: "}"}
	case ',':
		return token.Token{Type: token.Separator, Value: ","}
	case ':':
		return token.Token{Type: token.Separator, Value: ":"}
	case '"':
		return l.parseString()
	case 't':
		return l.parseKeyword(token.Bool, "true")
	case 'f':
		return l.parseKeyword(token.Bool, "false")
	case 'n':
		return l.parseKeyword(token.Null, "null")
	}
	if unicode.IsDigit(l.ch) || l.ch == '-' {
		return l.parseNumber()
	}
	if unicode.IsSpace(l.ch) {
		return l.parseWhitespace()
	}
	return token.Token{Type: token.UNKNOWN, Value: string(l.text[l.current])}
}
