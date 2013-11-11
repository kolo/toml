package toto

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	lexer    *lexer
	tree     map[string]interface{}
	keygroup string
}

func (p *parser) err(msg string) error {
	return &lexError{line: p.lexer.curLine, msg: msg}
}

func (p *parser) run() (err error) {
	var tok *token

	for {
		tok, err = p.lexer.nextToken()
		if err != nil {
			return err
		}

		if tok == nil {
			// EOF reached
			break
		}

		switch tok.tokenType {
		case tokKey:
			value, err := p.keyValue(tok.value)
			if err != nil {
				return err
			}
			p.setKey(tok.value, value)
		case tokKeyGroup:
			err = p.validateKeyGroup(tok.value)
			if err != nil {
				return err
			}
			p.keygroup = tok.value
		}
	}

	return nil
}

func (p *parser) setKey(k string, v interface{}) {
	if p.keygroup != "" {
		k = strings.Join([]string{p.keygroup, k}, ".")
	}
	p.tree[k] = v
}

func (p *parser) validateKeyGroup(keygroup string) error {
	keys := strings.Split(keygroup, ".")
	subkey := ""
	for _, k := range keys {
		if p.tree[k] != nil {
			return p.err("invalid keygroup")
		}
		if subkey != "" {
			subkey = strings.Join([]string{subkey, k}, ".")
		}
	}

	return nil
}

func (p *parser) keyValue(key string) (interface{}, error) {
	tok, err := p.lexer.nextToken()
	if err != nil {
		return nil, err
	}

	if tok == nil || tok.tokenType != tokAssignmentOperator {
		return nil, p.err("invalid key assignment")
	}
	p.lexer.assignment = true

	tok, err = p.lexer.nextToken()
	if err != nil {
		return nil, err
	}

	v, err := p.value(tok)
	if err != nil {
		return nil, err
	}
	p.lexer.assignment = false

	return v, nil
}

func (p *parser) arrayValue() ([]interface{}, error) {
	value := make([]interface{}, 0)

	t, err := p.lexer.nextToken()
	if err != nil {
		return nil, err
	}

	for {
		v, err := p.value(t)
		if err != nil {
			return nil, err
		}

		value = append(value, v)

		t, err = p.lexer.nextToken()
		if err != nil {
			return nil, err
		}

		if t.tokenType == tokRightBracket {
			break
		}

		if t.tokenType != tokComma {
			return nil, p.err("invalid array value")
		}

		t, err = p.lexer.nextToken()
		if err != nil {
			return nil, err
		}

		// There can be no last value. Example: a = [1,2,].
		if t.tokenType == tokRightBracket {
			break
		}
	}

	return value, nil
}

func (p *parser) value(t *token) (interface{}, error) {
	if t == nil || !isValueToken(t) {
		return nil, p.err("unknown value type")
	}

	switch t.tokenType {
	case tokLeftBracket:
		return p.arrayValue()
	case tokNumeric:
		return numericValue(t.value)
	case tokTrue:
		return true, nil
	case tokFalse:
		return false, nil
	default:
		return t.value, nil
	}
}

func isValueToken(t *token) bool {
	return t.tokenType == tokString || t.tokenType == tokNumeric ||
		t.tokenType == tokTrue || t.tokenType == tokFalse || t.tokenType == tokLeftBracket
}

var intValue = regexp.MustCompile(`^[0-9]+$`)
var floatValue = regexp.MustCompile(`^[0-9]+.[0-9]+$`)
var iso8601 = "2006-01-02T15:04:05Z"

func numericValue(v string) (interface{}, error) {
	if intValue.MatchString(v) {
		return strconv.ParseInt(v, 10, 64)
	}

	if floatValue.MatchString(v) {
		return strconv.ParseFloat(v, 64)
	}

	return time.Parse(iso8601, v)
}

func newParser(r io.Reader, tree map[string]interface{}) *parser {
	return &parser{
		lexer: newLexer(r),
		tree:  tree,
	}
}

func parse(r io.Reader, conf *Conf) (err error) {
	p := newParser(r, conf.values)
	err = p.run()
	return
}
