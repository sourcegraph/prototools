// Package util implements utilities for building protoc compiler plugins.
package util // import "sourcegraph.com/sourcegraph/prototools/util"

import (
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// ParseParams parses the comma-separated command-line parameters passed to the
// generator by protoc via r.GetParameters. Returned is a map of key=value
// parameters with whitespace preserved.
func ParseParams(r *plugin.CodeGeneratorRequest) map[string]string {
	// Split the parameter string and initialize the map.
	split := strings.Split(r.GetParameter(), ",")
	param := make(map[string]string, len(split))

	// Map the parameters.
	for _, p := range split {
		eq := strings.Split(p, "=")
		if len(eq) == 1 {
			param[strings.TrimSpace(eq[0])] = ""
			continue
		}
		val := strings.TrimSpace(eq[1])
		param[strings.TrimSpace(eq[0])] = val
	}
	return param
}
