package toml

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
enabled = true
verified = false

	[user.github]
	nickname = "TJ"
`

func Test_parseKeysExample(t *testing.T) {
	var i int64
	i = 5

	tests := parserTests{
		{"title", "TOML Example"},
		{"pool", i},
	}

	testParser(t, strings.NewReader(Keys), tests)
}

func Test_parseGroupsExample(t *testing.T) {
	tests := parserTests{
		{"title", "TOML Example"},
		{"user.name", "Tom Jones"},
		{"user.github.nickname", "TJ"},
		{"user.enabled", true},
		{"user.verified", false},
	}

	testParser(t, strings.NewReader(Groups), tests)
}

type parserTests []struct {
	key   string
	value interface{}
}

func testParser(t *testing.T, r io.Reader, tests parserTests) {
	keys := make(map[string]interface{})
	p := newParser(r, keys)
	err := p.run()

	if err != nil {
		t.Fatalf("run() returned unexpected error: %v", err)
	}

	for _, tt := range tests {
		if p.tree[tt.key] != tt.value {
			t.Errorf("%s is %v, though expected %v", tt.key, p.tree[tt.key], tt.value)
		}
	}
}
