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
	tokEOF
	tokComment
	tokKeyGroup
	tokKey
	tokString
	tokInt
	tokDate
	tokAssignmentOperator
	tokLeftBracket
	tokRightBracket
	tokSpace
)

type token struct {
	tokenType int
	value     string
}

type valueFunc func(rune) (string, error)

func selfValue(r rune) (string, error) {
	return string(r), nil
}

func newToken(tokenType int, r rune, value valueFunc) (t *token, err error) {
	t = new(token)
	t.tokenType = tokenType
	t.value, err = value(r)

	if err != nil {
		return nil, err
	}

	return
}

type lexer struct {
	scanner *bufio.Scanner
	atEOF   bool
}

func (l *lexer) nextToken() (t *token, err error) {
	r := l.next()
	if r == eof {
		return nil, nil
	}

	switch r {
	case '#':
		t, err = newToken(tokComment, r, l.commentValue)
	case '[':
		t, err = newToken(tokKeyGroup, r, l.keyGroupValue)
	case ' ', '\t':
		t, err = newToken(tokSpace, r, selfValue)
	case '=':
		t, err = newToken(tokAssignmentOperator, r, selfValue)
	case '"':
		t, err = newToken(tokString, r, l.stringValue)
	default:
		if unicode.IsLetter(r) {
			t, err = newToken(tokKey, r, l.keyValue)
		} else if unicode.IsDigit(r) {
			t, err = newToken(tokInt, r, l.intValue)
		} else {
			err = errors.New("unexpected token")
		}
	}

	if err != nil {
		return nil, err
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

func (l *lexer) scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	l.atEOF = atEOF
	advance, token, err = bufio.ScanRunes(data, atEOF)
	return
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

func (l *lexer) intValue(c rune) (string, error) {
	var buf bytes.Buffer
	buf.WriteRune(c)

	for {
		r := l.next()
		if r == ' ' || r == '\n' || r == eof {
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
		if r != ' ' && r != '\t' {
			return errors.New("unexpected character at the end of the line")
		}
	}
	return nil
}

func newLexer(r io.Reader) (l *lexer) {
	l = &lexer{
		atEOF: false,
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(l.scan)
	l.scanner = scanner

	return
}
