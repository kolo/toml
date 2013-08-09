package toto

import (
	"io"
)

type parser struct {
	lexer *lexer
}

func (p *parser) run() (err error) {
	var tok *token
	tok, err = p.lexer.nextToken()
	for err != nil {
		switch tok.tokenType {
		case tokComment:
			// Skip
		case tokEOF:
			break
		}

		tok, err = p.lexer.nextToken()
	}

	if err != nil {
		return err
	}

	return nil
}

func parse(r io.Reader, conf *Conf) (err error) {
	p := &parser{
		lexer: newLexer(r),
	}

	err = p.run()
	return
}
