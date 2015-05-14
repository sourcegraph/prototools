package util

import (
	"testing"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func TestParseParams(t *testing.T) {
	params := ParseParams(&plugin.CodeGeneratorRequest{
		Parameter: proto.String("key =value,abc = d ef , z = g "),
	})
	if len(params) != 3 {
		t.Fatal("expected 3 arguments got", len(params))
	}
	if params["key"] != "value" {
		t.Fatal(`"key" != "value"`)
	}
	if params["abc"] != "d ef" {
		t.Fatal(`"abc" != "d ef"`)
	}
	if params["z"] != "g" {
		t.Fatal(`"z" != "g"`)
	}
}

func TestIsFullyQualified(t *testing.T) {
	tests := map[string]bool{
		".google.protobuf.UninterpretedOption": true,
		".google.protobuf.FieldOptions.CType":  true,
		"protobuf.FieldOptions.CType":          false,
		"UninterpretedOption":                  false,
	}
	for symbolPath, want := range tests {
		got := IsFullyQualified(symbolPath)
		if got != want {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}
