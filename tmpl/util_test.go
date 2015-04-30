package tmpl

import "testing"

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
