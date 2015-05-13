# prototools

This repository holds various Protobuf/gRPC tools:

## [protoc-gen-doc](https://github.com/sourcegraph/prototools/blob/master/README.doc.md)

protoc-gen-doc is a Protobuf compiler plugin for generating HTML documentation from `.proto` files using standard Go HTML template files.

## [protoc-gen-json](https://github.com/sourcegraph/prototools/blob/master/README.json.md)

protoc-gen-json is a Protobuf compiler plugin for generating a JSON dump file for the `protoc` plugin generation request. It is primarily useful for running tests without invoking `protoc` directly (e.g. via the Go testing suite).
