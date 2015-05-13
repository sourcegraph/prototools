// Command protoc-gen-doc is a Protobuf plugin for documentation generation.
//
// Documentation can be found inside:
//
//  README.doc.md (https://github.com/sourcegraph/prototools/blob/master/README.doc.md)
//
// More information about Protobuf can be found at:
//
// 	https://developers.google.com/protocol-buffers/
//
package main // import "sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-doc"

import (
	"bytes"
	"go/build"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"sourcegraph.com/sourcegraph/prototools/tmpl"
	"sourcegraph.com/sourcegraph/prototools/util"
)

// PathDir returns the absolute path to a file given a relative one in one of
// the $GOPATH directories. If it cannot be resolved, relPath itself is
// returned.
func PathDir(relPath string) string {
	// Test again each directory listed in $GOPATH
	for _, path := range filepath.SplitList(build.Default.GOPATH) {
		path = filepath.Join(path, relPath)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return relPath
}

func main() {
	// Configure logging.
	log.SetFlags(0)
	log.SetPrefix("protoc-gen-doc: ")

	// Create a template generator.
	g := tmpl.New()

	// Read input from the protoc compiler.
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err, ": failed to read input")
	}

	// Unmarshal the protoc generation request.
	if err := proto.Unmarshal(data, g.Request); err != nil {
		log.Fatal(err, ": failed to parse input proto")
	}
	if len(g.Request.FileToGenerate) == 0 {
		log.Fatal(err, ": no input files")
	}

	// Parse the command-line parameters and determine the correct template file
	// to execute.
	params := util.ParseParams(g.Request)
	var tmplPath string
	if v, ok := params["template"]; ok {
		tmplPath = v
	} else {
		tmplPath = PathDir("src/sourcegraph.com/sourcegraph/prototools/templates/tmpl.html")
	}

	// Load up the template, preloading the function map (or else the functions
	// will fail when called).
	t := template.New("").Funcs(tmpl.Preload)
	g.Template, err = t.ParseGlob(tmplPath)
	if err != nil {
		log.Fatal(err, ": failed to parse templates")
	}
	g.Template = g.Template.Lookup(filepath.Base(tmplPath))

	// Perform generation.
	response, err := g.Generate()
	if err != nil {
		log.Fatal(err, ": failed to generate")
	}

	// Marshal the results and write back to the protoc compiler.
	data, err = proto.Marshal(response)
	if err != nil {
		log.Fatal(err, ": failed to marshal output proto")
	}
	_, err = io.Copy(os.Stdout, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err, ": failed to write output proto")
	}
}
