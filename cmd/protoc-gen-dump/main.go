// Command protoc-gen-dump is a Protobuf plugin for generating proto dump files.
//
// Documentation can be found inside:
//
//  README.dump.md (https://github.com/sourcegraph/prototools/blob/master/README.dump.md)
//
// More information about Protobuf can be found at:
//
// 	https://developers.google.com/protocol-buffers/
//
package main // import "sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-dump"

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"sourcegraph.com/sourcegraph/prototools/util"
)

func main() {
	// Configure logging.
	log.SetFlags(0)
	log.SetPrefix("protoc-gen-proto: ")

	// Read input from the protoc compiler.
	data, err := io.ReadAll(os.Stdin)
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
	name := "out.gob"
	if v, ok := params["out"]; ok {
		name = v
	}

	// Perform generation.
	response := &plugin.CodeGeneratorResponse{}
	data, err = proto.Marshal(request)
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
