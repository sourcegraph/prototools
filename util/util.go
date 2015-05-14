// Package util implements utilities for building protoc compiler plugins.
package util // import "sourcegraph.com/sourcegraph/prototools/util"

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

// FieldTypeName returns the protobuf-syntax name for the given field type. It
// panics on errors (e.g. zero value).
func FieldTypeName(f *descriptor.FieldDescriptorProto_Type) string {
	switch *f {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "double"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "group"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return "message"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "bytes"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "enum"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "sfixed32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "sfixed64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "sint32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "sint64"
	default:
		panic("FieldTypeName: unknown field type")
	}
}

// IsFullyQualified tells if the given symbol path is fully-qualified or not (i.e.
// starts with a period).
func IsFullyQualified(symbolPath string) bool {
	return symbolPath[0] == '.'
}
