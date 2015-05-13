// Command protoc-gen-doc is a Protobuf plugin for generating JSON dump files.
//
// Documentation can be found inside:
//
//  README.json.md (https://github.com/sourcegraph/prototools/blob/master/README.json.md)
//
// More information about Protobuf can be found at:
//
// 	https://developers.google.com/protocol-buffers/
//
package main // import "sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-json"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"sourcegraph.com/sourcegraph/prototools/util"
)

func main() {
	// Configure logging.
	log.SetFlags(0)
	log.SetPrefix("protoc-gen-json: ")

	// Read input from the protoc compiler.
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err, ": failed to read input")
	}

	// Unmarshal the protoc generation request.
	request := &plugin.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, request); err != nil {
		log.Fatal(err, ": failed to parse input proto")
	}
	if len(request.FileToGenerate) == 0 {
		log.Fatal(err, ": no input files")
	}

	// Parse the command-line parameters.
	params := util.ParseParams(request)

	// Determine the correct output file name.
	name := "out.json"
	if v, ok := params["out"]; ok {
		name = v
	}

	// Determine the JSON indention.
	indent := "  "
	if v, ok := params["indent"]; ok {
		indent = v
	}

	// Perform generation.
	response := &plugin.CodeGeneratorResponse{}
	data, err = json.MarshalIndent(request, "", indent)
	if err != nil {
		response.Error = proto.String(err.Error())
	} else {
		response.File = []*plugin.CodeGeneratorResponse_File{
			&plugin.CodeGeneratorResponse_File{
				Name:    proto.String(name),
				Content: proto.String(string(data)),
			},
		}
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
