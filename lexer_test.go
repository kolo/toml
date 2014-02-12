package toml

import (
	"strings"
	"testing"
)

func Test_readKeyGroup(t *testing.T) {
	testToken(t, "[key.subkey]", tokKeyGroup, "key.subkey")
}

func Test_readKey(t *testing.T) {
	testToken(t, "key = ", tokKey, "key")
}

func Test_readString(t *testing.T) {
	testToken(t, "\"I'm a string. \\\"You can quote me\\\". Name\\tJos\\u00E9\\nLocation\\tSF.\"",
		tokString, "I'm a string. \\\"You can quote me\\\". Name\\tJos\\u00E9\\nLocation\\tSF.")
}

func Test_readInt(t *testing.T) {
	testToken(t, "42 ", tokNumeric, "42")
	testToken(t, "42.5 \n", tokNumeric, "42.5")
}

func testToken(t *testing.T, src string, tokenType int, expected string) {
	r := strings.NewReader(src)
	l := newLexer(r)

	tt, err := l.nextToken()
	if err != nil {
		t.Fatalf("nextToken() returned unexpected error: %v", err)
	}
	if tt == nil {
		t.Fatal("nextToken() returned nil")
	}
	if tt.tokenType != tokenType {
		t.Fatalf("Expected %d, but got %d", tokenType, tt.tokenType)
	}
	if tt.value != expected {
		t.Fatalf("Expect key to be \"%s\", but got \"%s\"", expected, tt.value)
	}
}
