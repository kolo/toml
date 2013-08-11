package toto

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"unicode"
	"unicode/utf8"
)

const eof = -(iota + 1)

const (
	tokError = iota
	tokKeyGroup
	tokKey
	tokString
	tokNumeric
	tokTrue
	tokFalse
	tokArray
	tokAssignmentOperator
)

type token struct {
	tokenType int
	value     string
}

type valueFunc func(rune) (string, error)

func selfValue(r rune) (string, error) {
	return string(r), nil
}

func trueValue(rune) (string, error) {
	return "true", nil
}

func falseValue(rune) (string, error) {
	return "false", nil
}

type lexer struct {
	scanner   *bufio.Scanner
	lastToken *token
}

func (l *lexer) nextToken() (t *token, err error) {
	for {
		r := l.next()
		if r == eof {
			return nil, nil
		}

		switch r {
		case ' ', '\t', '\n':
			// Skip
		case '#':
			l.omitLineReminder()
		case '[':
			t, err = l.newToken(tokKeyGroup, r, l.keyGroupValue)
		case '=':
			t, err = l.newToken(tokAssignmentOperator, r, selfValue)
		case '"':
			t, err = l.newToken(tokString, r, l.stringValue)
		default:
			if unicode.IsLetter(r) {
				if l.lastToken != nil && l.lastToken.tokenType == tokAssignmentOperator {
					t, err = l.newBooleanToken(r)
				} else {
					t, err = l.newToken(tokKey, r, l.keyValue)
				}
			} else if unicode.IsDigit(r) {
				t, err = l.newToken(tokNumeric, r, l.value)
			} else {
				err = errors.New("unexpected token")
			}
		}

		if err != nil || t != nil {
			break
		}
	}

	return
}

func (l *lexer) next() rune {
	ok := l.scanner.Scan()
	if !ok {
		return eof
	}
	r, _ := utf8.DecodeRune(l.scanner.Bytes())

	return r
}

func (l *lexer) newToken(tokenType int, r rune, value valueFunc) (t *token, err error) {
	t = new(token)
	t.tokenType = tokenType
	t.value, err = value(r)

	if err != nil {
		return nil, err
	}

	l.lastToken = t

	return
}

func (l *lexer) newBooleanToken(r rune) (t *token, err error) {
	v, err := l.value(r)
	if err != nil {
		return nil, err
	}

	switch v {
	case "true":
		return l.newToken(tokTrue, r, trueValue)
	case "false":
		return l.newToken(tokFalse, r, falseValue)
	default:
		return nil, errors.New("unknown value type")
	}
}

func (l *lexer) commentValue(rune) (string, error) {
	var buf bytes.Buffer

	for {
		r := l.next()
		if r == '\n' || r == eof {
			break
		}
		buf.WriteRune(r)
	}

	return buf.String(), nil
}

func (l *lexer) keyGroupValue(rune) (string, error) {
	var buf bytes.Buffer
	var newKey, finished bool

	newKey = true
	finished = false

	for !finished {
		r := l.next()
		if newKey {
			if !unicode.IsLetter(r) {
				return "", errors.New("invalid keygroup")
			}
			newKey = false
			buf.WriteRune(r)
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				buf.WriteRune(r)
			}
			if r == ']' {
				finished = true
			}
			if r == '.' {
				newKey = true
				buf.WriteRune(r)
			}
		}
	}

	err := l.omitLineReminder()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (l *lexer) keyValue(c rune) (string, error) {
	var buf bytes.Buffer
	var finished bool

	finished = false
	buf.WriteRune(c)

	for !finished {
		r := l.next()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			buf.WriteRune(r)
		}
		if r == ' ' {
			finished = true
		}
	}

	return buf.String(), nil
}

func (l *lexer) stringValue(rune) (string, error) {
	var buf bytes.Buffer
	var err error

	escaped := false

	for {
		r := l.next()
		if !escaped && r == '"' {
			break
		}
		if r == '\n' || r == eof {
			err = errors.New("unexpected end of line")
			break
		}
		if escaped {
			if r == 'b' || r == 't' || r == 'n' || r == 'f' || r == 'r' ||
				r == '"' || r == '/' || r == '\\' || r == 'u' {
				escaped = false
				buf.WriteRune(r)
			} else {
				err = errors.New("unknown escape sequence")
				break
			}
		} else {
			if r == '\\' {
				escaped = true
			}
			buf.WriteRune(r)
		}
	}

	if err != nil {
		return "", err
	}

	err = l.omitLineReminder()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (l *lexer) value(c rune) (string, error) {
	var buf bytes.Buffer
	buf.WriteRune(c)

	for {
		r := l.next()
		if isSpace(r) || r == '\n' || r == eof {
			break
		}
		buf.WriteRune(r)
	}

	return buf.String(), nil
}

func (l *lexer) omitLineReminder() error {
	// Only space and tab allowed at the reminder of the line.
	for {
		r := l.next()
		if r == '\n' || r == eof {
			break
		}
		if !isSpace(r) {
			return errors.New("unexpected character at the end of the line")
		}
	}
	return nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func newLexer(r io.Reader) (l *lexer) {
	l = &lexer{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	l.scanner = scanner

	return
}
