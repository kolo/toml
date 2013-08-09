package toto

import (
	"io"
	"strings"
	"testing"
)

const ValuesOnly = `
title = "TOML Example"

pool = 5
`

func Test_parseExample(t *testing.T) {
	tests := parserTests{
		{"title", "TOML Example"},
		{"pool", "5"},
	}

	testParser(t, strings.NewReader(ValuesOnly), tests)
}

type parserTests []struct{
	key string
	value string
}

func testParser(t *testing.T, r io.Reader, tests parserTests) {
	p := newParser(r)
	err := p.run()

	if err != nil {
		t.Fatalf("run() returned unexpected error: %v", err)
	}


	for _, tt := range tests {
		if p.tree[tt.key] != tt.value {
			t.Errorf("%s is %s, though expected %s", tt.key, p.tree[tt.key], tt.value)
		}
	}
}
