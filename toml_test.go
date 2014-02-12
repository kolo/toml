package toml

import (
	"testing"
)

func Test_valueMethods(t *testing.T) {
	conf, err := Parse("example.toml")
	if err != nil {
		t.Fatal("Parse() returned unexpected error: %v", err)
	}

	var connection_max int64
	connection_max = 5000

	results := map[string]interface{}{
		"title": "TOML Example",
		"database.connection_max": connection_max,
		"database.enabled": true,
		"clients.hosts": []string{"alpha", "beta"},
	}

	if v := conf.String("title"); v != results["title"] {
		t.Fatalf("title is %s, though expected %v", v, results["title"])
	}

	if v := conf.Int("database.connection_max"); v != results["database.connection_max"] {
		t.Fatalf("database.connection_max is %v, though expected %v", v, results["database.connection_max"])
	}

	if v := conf.Bool("database.enabled"); v != results["database.enabled"] {
		t.Fatalf("database.enabled is %v, though expected %v", v, results["database.enabled"])
	}

	v := conf.Slice("clients.hosts")
	if v == nil {
		t.Fatal("clients.hosts is nil, though expected %v", results["clients.hosts"])
	}
}

func Test_accessToUndefinedKeys(t *testing.T) {
	conf, err := Parse("example.toml")
	if err != nil {
		t.Fatal("Parse() returned unexpected error: %v", err)
	}

	var emptyString string
	if v := conf.String("undefined.string"); v != emptyString {
		t.Fatal("undefined.string key should be empty")
	}

	var zero int64
	if v := conf.Int("undefined.int"); v != zero {
		t.Fatal("undefined.int key should be empty")
	}

	var boolean bool
	if v := conf.Bool("undefined.bool"); v != boolean {
		t.Fatal("undefined.bool key should be empty")
	}

	if v := conf.Slice("undefined.slice"); v != nil {
		t.Fatal("undefined.slice key should be nil")
	}
}
