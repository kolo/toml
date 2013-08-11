package toto

import (
	"errors"
	"io"
	"strings"
)

type parser struct {
	lexer    *lexer
	tree     map[string]interface{}
	keygroup string
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
			err = p.setKeyGroup(tok.value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *parser) setKey(k string, v string) {
	if p.keygroup != "" {
		k = strings.Join([]string{p.keygroup, k}, ".")
	}
	p.tree[k] = v
}

func (p *parser) setKeyGroup(keygroup string) error {
	if p.keygroup != "" {
		keygroup = strings.Join([]string{p.keygroup, keygroup}, ".")
	}

	keys := strings.Split(keygroup, ".")
	subkey := ""
	for _, k := range keys {
		if p.tree[k] != nil {
			return errors.New("invalid keygroup")
		}
		if subkey != "" {
			subkey = strings.Join([]string{subkey, k}, ".")
		}
	}

	p.keygroup = keygroup

	return nil
}

func (p *parser) keyValue(key string) (string, error) {
	tok, err := p.lexer.nextToken()
	if err != nil {
		return "", err
	}
	if tok == nil || tok.tokenType != tokAssignmentOperator {
		return "", errors.New("invalid key assignment")
	}

	tok, err = p.lexer.nextToken()
	if err != nil {
		return "", err
	}
	if tok == nil || !isValueToken(tok) {
		return "", errors.New("invalid key assignment")
	}

	return tok.value, nil
}

func isValueToken(t *token) bool {
	return t.tokenType == tokString || t.tokenType == tokInt ||
		t.tokenType == tokFloat || t.tokenType == tokDate
}

func newParser(r io.Reader) *parser {
	return &parser{
		lexer: newLexer(r),
		tree:  make(map[string]interface{}),
	}
}

func parse(r io.Reader, conf *Conf) (err error) {
	p := newParser(r)
	err = p.run()
	return
}
