package toto

import (
	"errors"
	"io"
)

type parser struct {
	lexer *lexer
	tree  map[string]interface{}
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
			p.tree[tok.value] = value
		}
	}

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
		tree: make(map[string]interface{}),
	}
}

func parse(r io.Reader, conf *Conf) (err error) {
	p := newParser(r)
	err = p.run()
	return
}
