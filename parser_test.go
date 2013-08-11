package toto

import (
	"io"
	"strings"
	"testing"
)

const Keys = `
title = "TOML Example"

pool = 5
`

const Groups = `
title = "TOML Example"

[user]
name = "Tom Jones"

	[github]
	nickname = "TJ"
`

func Test_parseKeysExample(t *testing.T) {
	tests := parserTests{
		{"title", "TOML Example"},
		{"pool", "5"},
	}

	testParser(t, strings.NewReader(Keys), tests)
}

func Test_parseGroupsExample(t *testing.T) {
	tests := parserTests {
		{"title", "TOML Example"},
		{"user.name", "Tom Jones"},
		{"user.github.nickname", "TJ"},
	}

	testParser(t, strings.NewReader(Groups), tests)
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
			t.Errorf("%s is %v, though expected %s", tt.key, p.tree[tt.key], tt.value)
		}
	}
}
