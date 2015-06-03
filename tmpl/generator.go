// Package tmpl implements a protobuf template-based generator.
package tmpl // import "sourcegraph.com/sourcegraph/prototools/tmpl"

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"path"

	gateway "github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generator is the type whose methods generate the output, stored in the associated response structure.
type Generator struct {
	// Request from protoc compiler, which should have data unmarshaled into it by
	// the user of this package.
	Request *plugin.CodeGeneratorRequest

	// FileMap is the map of template files to use for the generation process.
	FileMap FileMap

	// RootDir is the root directory path prefix to place onto URLs for generated
	// types.
	RootDir string

	// APIHost is the base URL to use for rendering grpc-gateway routes, e.g.:
	//
	//  http://api.mysite.com/
	//
	APIHost string

	// Response to protoc compiler.
	response *plugin.CodeGeneratorResponse
}

// ParseFileMap parses and executes a filemap template.
func (g *Generator) ParseFileMap(dir, data string) error {
	// Parse the template data.
	t, err := template.New("").Funcs(Preload).Parse(data)
	if err != nil {
		return err
	}

	// Execute the template.
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, g.Request)
	if err != nil {
		return err
	}

	// Parse the filemap.
	g.FileMap.Dir = dir
	err = xml.Unmarshal(buf.Bytes(), &g.FileMap)
	if err != nil {
		return err
	}
	if len(g.FileMap.Generate) == 0 {
		return errors.New("no generate elements found in file map")
	}
	return nil
}

// Generate generates a response for g.Request (which you should unmarshal data
// into using protobuf).
//
// If any error is encountered during generation, it is returned and should be
// considered fatal to the generation process (the response will be nil).
func (g *Generator) Generate() (response *plugin.CodeGeneratorResponse, err error) {
	// Reset the response to its initial state.
	g.response.Reset()

	// Create a grpc-gateway registry.
	reg := gateway.NewRegistry()
	if err := reg.Load(g.Request); err != nil {
		return nil, err
	}

	// Generate each proto file:
	errs := new(bytes.Buffer)
	buf := new(bytes.Buffer)
	protoFile := g.Request.GetProtoFile()
	for _, f := range protoFile {
		for _, gen := range g.FileMap.Generate {
			// Only running generators on proto files (i.e. generators with
			// targets).
			if gen.Target != f.GetName() {
				continue
			}

			// Prepare the generators template.
			tmpl, err := g.prepare(gen)
			if err != nil {
				fmt.Fprintf(errs, "%s\n", err)
				continue
			}

			// Execute the template with this context and generate a response
			// for the input file.
			buf.Reset()
			ctx := &tmplFuncs{
				f:          f,
				outputFile: gen.Output,
				rootDir:    g.RootDir,
				protoFile:  protoFile,
				registry:   reg,
				apiHost:    g.APIHost,
			}
			err = tmpl.Funcs(ctx.funcMap()).Execute(buf, struct {
				*descriptor.FileDescriptorProto
				Generate *FileMapGenerate
				Data     map[string]string
				Request  *plugin.CodeGeneratorRequest
			}{
				f,
				gen,
				gen.DataMap(),
				g.Request,
			})
			if err != nil {
				fmt.Fprintf(errs, "%s\n", err)
				continue
			}

			// Generate the response file with the rendered template.
			bufStr := buf.String()
			g.response.File = append(g.response.File, &plugin.CodeGeneratorResponse_File{
				Name:    &gen.Output,
				Content: &bufStr,
			})
		}
	}

	// Execute target-less filemap generators (e.g. for index pages rather than
	// individual doc pages).
	for _, gen := range g.FileMap.Generate {
		// Only running generators not on proto files (i.e. generators without
		// targets).
		if len(gen.Target) != 0 {
			continue
		}

		// Prepare the generators template.
		tmpl, err := g.prepare(gen)
		if err != nil {
			fmt.Fprintf(errs, "%s\n", err)
			continue
		}

		// Execute the template with this context and generate a response file.
		buf.Reset()
		ctx := &tmplFuncs{
			outputFile: gen.Output,
			rootDir:    g.RootDir,
			registry:   reg,
			apiHost:    g.APIHost,
		}
		err = tmpl.Funcs(ctx.funcMap()).Execute(buf, struct {
			*plugin.CodeGeneratorRequest
			Generate *FileMapGenerate
			Data     map[string]string
		}{
			g.Request,
			gen,
			gen.DataMap(),
		})
		if err != nil {
			fmt.Fprintf(errs, "%s\n", err)
			continue
		}

		// Generate the response file with the rendered template.
		bufStr := buf.String()
		g.response.File = append(g.response.File, &plugin.CodeGeneratorResponse_File{
			Name:    &gen.Output,
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

// New returns a new generator for the given template.
func New() *Generator {
	return &Generator{
		Request:  &plugin.CodeGeneratorRequest{},
		response: &plugin.CodeGeneratorResponse{},
	}
}

// prepare prepares the given filemap generators template for execution,
// handling parsing of both the relative-path templates and their includes.
func (g *Generator) prepare(gen *FileMapGenerate) (*template.Template, error) {
	// Preload the function map (or else the functions will fail when
	// called due to a lack of valid context).
	var (
		t   = template.New("").Funcs(Preload)
		err error
	)

	// Parse the included template files.
	if len(gen.Include) > 0 {
		t, err = t.ParseFiles(g.FileMap.relative(gen.Include...)...)
		if err != nil {
			return nil, err
		}
	}

	// Parse the template file to execute.
	absTemplate := g.FileMap.relative(gen.Template)[0]
	tmpl, err := t.ParseFiles(absTemplate)
	if err != nil {
		return nil, err
	}
	_, name := path.Split(absTemplate)
	tmpl = tmpl.Lookup(name)
	return tmpl, nil
}
