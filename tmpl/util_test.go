package tmpl

import (
	"bytes"
	"testing"
)

func TestUnixPath(t *testing.T) {
	var paths = map[string]string{
		"a\\b\\c\\d/e\\": "a/b/c/d/e",
		"a/b/c/":         "a/b/c",
		"a/b":            "a/b",
	}
	for p, correct := range paths {
		got := unixPath(p)
		if got != correct {
			t.Fatalf("got %q expected %q", got, correct)
		}
	}
}

func TestStripExt(t *testing.T) {
	var files = map[string]string{
		"hello.txt":      "hello",
		"no_ext":         "no_ext",
		"long.extension": "long",
	}
	for f, correct := range files {
		got := stripExt(f)
		if got != correct {
			t.Fatalf("got %q expected %q", got, correct)
		}
	}
}

func TestPkgStmt(t *testing.T) {
	var tests = [][2]string{
		{"package foobar3;\n", "foobar3"},
		{"package foobar3   \t; \n", "foobar3"},
		{"package FooBar3;\n", "FooBar3"},
		{"package foo-bar3;\n", ""},
		{"package foo3; // comments are silly", "foo3"},
		{" \t package foobar3;\n", "foobar3"},
		{"BAH\n \t package foo4;\n", "foo4"},
		{"// package notIt;\n  package 3eggs3;\n", "3eggs3"},
	}
	for _, tst := range tests {
		pkg, err := pkgStmt(bytes.NewReader([]byte(tst[0])))
		if err != nil {
			t.Fatal(err)
		}
		wantPkg := tst[1]
		if pkg != wantPkg {
			t.Logf("%q\n", tst[0])
			t.Fatalf("got=%q want=%q\n", pkg, wantPkg)
		}
	}
}
