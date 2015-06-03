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
	"encoding/xml"
	"fmt"
	"go/build"
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

// extendParams extends the given parameter map with the second one.
func extendParams(params, second map[string]string) map[string]string {
	for k, v := range second {
		if _, ok := params[k]; !ok {
			params[k] = v
		}
	}
	return params
}

var basicFileMap = `
<FileMap>
{{$templatePath := "%s"}}
{{range .ProtoFile}}
    <Generate>
        <Template>{{$templatePath}}</Template>
        <Target>{{.Name}}</Target>
        <Output>{{trimExt .Name}}{{ext $templatePath}}</Output>
    </Generate>
{{end}}
</FileMap>
`

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

	// Verify the command-line parameters.
	params := util.ParseParams(g.Request)

	// Handle configuration files.
	if conf, ok := params["conf"]; ok {
		confData, err := ioutil.ReadFile(conf)
		if err != nil {
			log.Fatal(err, ": could not read conf file")
		}
		g.Request.Parameter = proto.String(string(confData))
		params = extendParams(params, util.ParseParams(g.Request))
	}

	paramTemplate, haveTemplate := params["template"]
	paramFileMap, haveFileMap := params["filemap"]
	if haveTemplate && haveFileMap {
		log.Fatal("expected either template or filemap argument, not both")
	}

	// Build the filemap based on the command-line parameters.
	var fileMapDir, fileMapData string
	if haveTemplate {
		// Use the specified template file once on each input proto file.
		fileMapData = fmt.Sprintf(basicFileMap, paramTemplate)
	} else if haveFileMap {
		// Load the filemap template.
		data, err := ioutil.ReadFile(paramFileMap)
		if err != nil {
			log.Fatal(err, ": failed to read file map")
		}
		fileMapData = string(data)
		fileMapDir = filepath.Dir(paramFileMap)
	} else {
		// Use the default filemap template once on each input proto file.
		def := PathDir("src/sourcegraph.com/sourcegraph/prototools/templates/tmpl.html")
		fileMapData = fmt.Sprintf(basicFileMap, def)
		fileMapDir = filepath.Dir(def)
	}

	// Parse the file map template.
	if err = g.ParseFileMap(fileMapDir, fileMapData); err != nil {
		log.Fatal(err, ": failed to parse file map")
	}

	// Dump the execute filemap template, if desired.
	if v, ok := params["dump-filemap"]; ok {
		f, err := os.Create(v)
		if err != nil {
			log.Fatal(err, ": failed to crate dump file")
		}
		dump, err := xml.MarshalIndent(g.FileMap, "", "    ")
		if err != nil {
			log.Fatal(err, ": failed to marshal filemap")
		}
		_, err = io.Copy(f, bytes.NewReader(dump))
		if err != nil {
			log.Fatal(err, ": failed to write dump file")
		}
	}

	// Determine the root directory.
	if v, ok := params["root"]; ok {
		g.RootDir = v
	} else {
		g.RootDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Map the API host, if any.
	if v, ok := params["apihost"]; ok {
		g.APIHost = v
	}

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
