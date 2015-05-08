// Package tmpl implements a protobuf template-based generator.
package tmpl // import "sourcegraph.com/sourcegraph/prototools/tmpl"

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generator is the type whose methods generate the output, stored in the associated response structure.
type Generator struct {
	// Request from protoc compiler, which should have data unmarshaled into it by
	// the user of this package.
	Request *plugin.CodeGeneratorRequest

	// The template to execute for generation, which should be set by the user of
	// this package.
	Template *template.Template

	// The command-line parameters passed to the generator via protoc; initialized
	// after g.ParseParams or g.Generate has been called.
	Param map[string]string

	// Extension is the extension string to name generated files with. If it is an
	// empty string, the extension of the first template file is used.
	Extension string

	// Response to protoc compiler.
	response *plugin.CodeGeneratorResponse
}

// Generate generates a response for g.Request (which you should unmarshal data
// into using protobuf).
//
// If any error is encountered during generation, it is returned and should be
// considered fatal to the generation process (the response will be nil).
func (g *Generator) Generate() (response *plugin.CodeGeneratorResponse, err error) {
	// Reset the response to it's initial state.
	g.response.Reset()

	// Parse command-line parameters.
	g.ParseParams()

	// Determine the extension string.
	ext := g.Extension
	if len(ext) == 0 {
		ext = findExt(g.Template)
	}

	// Generate each proto file:
	errs := new(bytes.Buffer)
	buf := new(bytes.Buffer)
	for _, f := range g.Request.GetProtoFile() {
		// Execute the template and generate a response for the input file.
		buf.Reset()
		err := g.Template.Funcs(newTmplFuncs(f, ext).funcMap()).Execute(buf, f)

		// If an error occured during executing the template, we pass it pack to
		// protoc via the error field in the response.
		if err != nil {
			fmt.Fprintf(errs, "%s\n", err)
			continue
		}

		// Determine the file name (relative to the output directory).
		name := stripExt(f.GetName()) + ext
		name = unixPath(name)

		// Generate the response file with the rendered template.
		bufStr := buf.String()
		g.response.File = append(g.response.File, &plugin.CodeGeneratorResponse_File{
			Name:    &name,
			Content: &bufStr,
		})
	}
	if errs.Len() > 0 {
		g.response.File = nil
		errsStr := errs.String()
		g.response.Error = &errsStr
	}
	return g.response, nil
}

// ParseParams parses the command-line parameters passed to the generator by
// protoc via g.Request.GetParameters. It can be called as soon as g.Request is
// assigned; and is automatically called at generation time.
func (g *Generator) ParseParams() {
	// Split the parameter string and initialize the map.
	split := strings.Split(g.Request.GetParameter(), ",")
	g.Param = make(map[string]string, len(split))

	// Map the parameters.
	for _, p := range split {
		if i := strings.Index(p, "="); i < 0 {
			g.Param[p] = ""
		} else {
			g.Param[p[0:i]] = p[i+1:]
		}
	}
}

// New returns a new generator for the given template.
func New() *Generator {
	return &Generator{
		Request:  &plugin.CodeGeneratorRequest{},
		response: &plugin.CodeGeneratorResponse{},
	}
}
