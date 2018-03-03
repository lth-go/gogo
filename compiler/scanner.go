package compiler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const (
	// EOF End of file.
	EOF = -1
	// EOL End of line.
	EOL = '\n'
)

// opName is correction of operation names.
var opName = map[string]int{
	"if":       IF,
	"else":     ELSE,
	"elif":     ELIF,
	"for":      FOR,
	"return":   RETURN_T,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE_T,
	"false":    FALSE_T,
	"void":     VOID_T,
	"boolean":  BOOLEAN_T,
	"int":      INT_T,
	"double":   DOUBLE_T,
	"string":   STRING_T,
	"null":     NULL_T,
	"new":      NEW,
	"require":  REQUIRE,
	"class":    CLASS_T,
	"this":     THIS_T,
	"(":        LP,
	")":        RP,
	"[":        LB,
	"]":        RB,
	"{":        LC,
	"}":        RC,
	";":        SEMICOLON,
	":":        COLON,
	",":        COMMA,
	"+":        ADD,
	"-":        SUB,
	"*":        MUL,
	"/":        DIV,
	"!":        EXCLAMATION,
	".":        DOT,
}

// Scanner stores informations for lexer.
type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func newScannerByFilePath(path string) *Scanner {
	_, err := os.Stat(path)
	if err != nil {
		panic("文件不存在")
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	scanner := &Scanner{src: []rune(string(buf))}

	return scanner
}

// Scan analyses token, and decide identify or literals.
func (s *Scanner) Scan() (tok int, lit string, pos Position, err error) {
retry:
	s.skipBlank()
	pos = s.pos()
	switch ch := s.peek(); {
	// 关键字
	case isLetter(ch):
		lit, err = s.scanIdentifier()
		if err != nil {
			return
		}
		if name, ok := opName[lit]; ok {
			tok = name
		} else {
			tok = IDENTIFIER
		}
	// 数字
	case isDigit(ch):
		lit, err = s.scanNumber()
		if err != nil {
			return
		}
		// 判断lit中有无 `.`
		if strings.Contains(lit, ".") {
			tok = DOUBLE_LITERAL
		} else {
			tok = INT_LITERAL
		}
	// 字符串
	case ch == '"':
		tok = STRING_LITERAL
		lit, err = s.scanString('"')
		if err != nil {
			return
		}
	case ch == '\'':
		tok = STRING_LITERAL
		lit, err = s.scanString('\'')
		if err != nil {
			return
		}
	default:
		switch ch {
		case EOF:
			tok = EOF
		case '\n':
			s.next()
			goto retry
		// 注释
		case '#':
			for !isEOL(s.peek()) {
				s.next()
			}
			goto retry
		case '=':
			s.next()
			switch s.peek() {
			case '=':
				tok = EQ
				lit = "=="
			default:
				s.back()
				tok = ASSIGN_T
				lit = "="
			}
		case '!':
			s.next()
			switch s.peek() {
			case '=':
				tok = NE
				lit = "!="
			default:
				s.back()
				tok = EXCLAMATION
				lit = "!"
			}
		case '>':
			s.next()
			switch s.peek() {
			case '=':
				tok = GE
				lit = ">="
			default:
				s.back()
				tok = GT
				lit = ">"
			}
		case '<':
			s.next()
			switch s.peek() {
			case '=':
				tok = LE
				lit = "<="
			default:
				s.back()
				tok = LT
				lit = "<"
			}
		case '|':
			s.next()
			switch s.peek() {
			case '|':
				tok = LOGICAL_OR
				lit = "||"
			default:
				s.back()
				err = fmt.Errorf(`Syntax Error "%s"`, string(ch))
				tok = int(ch)
				lit = string(ch)
			}
		case '&':
			s.next()
			switch s.peek() {
			case '&':
				tok = LOGICAL_AND
				lit = "&&"
			default:
				s.back()
				err = fmt.Errorf(`Syntax Error "%s"`, string(ch))
				tok = int(ch)
				lit = string(ch)
			}
		case '(', ')', '[', ']', '{', '}', ':', ';', ',', '+', '-', '*', '/', '.':
			tok = opName[string(ch)]
			lit = string(ch)
		default:
			err = fmt.Errorf(`Syntax Error "%s"`, string(ch))
			tok = int(ch)
			lit = string(ch)
			return
		}
		s.next()
	}
	return
}

// ==============================
// isXxx
// ==============================

// isLetter returns true if the rune is a letter for identity.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isDigit returns true if the rune is a number.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isEOL returns true if the rune is at end-of-line or end-of-file.
func isEOL(ch rune) bool {
	return ch == '\n' || ch == -1
}

// isBlank returns true if the rune is empty character..
func isBlank(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// ==============================
// move
// ==============================

// peek returns current rune in the code.
func (s *Scanner) peek() rune {
	if s.reachEOF() {
		return EOF
	}
	return s.src[s.offset]
}

// next moves offset to next.
func (s *Scanner) next() {
	if s.reachEOF() {
		return
	}

	if s.peek() == '\n' {
		s.lineHead = s.offset + 1
		s.line++
	}
	s.offset++
}

// back moves back offset once to top.
func (s *Scanner) back() {
	s.offset--
}

// skipBlank moves position into non-black character.
func (s *Scanner) skipBlank() {
	for isBlank(s.peek()) {
		s.next()
	}
}

// reachEOF returns true if offset is at end-of-file.
func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

// pos returns the position of current.
func (s *Scanner) pos() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

// ==============================
// scanXXX
// ==============================

// scanIdentifier returns identifier begining at current position.
func (s *Scanner) scanIdentifier() (string, error) {
	var ret []rune
	for {
		if !isLetter(s.peek()) && !isDigit(s.peek()) {
			break
		}
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret), nil
}

// scanNumber returns number begining at current position.
func (s *Scanner) scanNumber() (string, error) {
	var ret []rune
	ch := s.peek()
	ret = append(ret, ch)
	s.next()

	for isDigit(s.peek()) || s.peek() == '.' {
		ret = append(ret, s.peek())
		s.next()
	}

	if isLetter(s.peek()) {
		return "", errors.New("identifier starts immediately after numeric literal")
	}
	return string(ret), nil
}

// scanString returns string starting at current position.
// This handles backslash escaping.
func (s *Scanner) scanString(l rune) (string, error) {
	var ret []rune
eos:
	for {
		s.next()
		switch s.peek() {
		case EOL:
			return "", errors.New("unexpected EOL")
		case EOF:
			return "", errors.New("unexpected EOF")
		case l:
			s.next()
			break eos
		case '\\':
			s.next()
			switch s.peek() {
			case 'b':
				ret = append(ret, '\b')
				continue
			case 'f':
				ret = append(ret, '\f')
				continue
			case 'r':
				ret = append(ret, '\r')
				continue
			case 'n':
				ret = append(ret, '\n')
				continue
			case 't':
				ret = append(ret, '\t')
				continue
			}
			ret = append(ret, s.peek())
			continue
		default:
			ret = append(ret, s.peek())
		}
	}
	return string(ret), nil
}
