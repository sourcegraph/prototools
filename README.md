# prototools

This repository holds various Protobuf/gRPC tools:

## [protoc-gen-doc](https://github.com/sourcegraph/prototools/blob/master/README.doc.md)

protoc-gen-doc is a Protobuf compiler plugin for generating HTML documentation from `.proto` files using standard Go HTML template files.

## [protoc-gen-json](https://github.com/sourcegraph/prototools/blob/master/README.json.md)

protoc-gen-json is a Protobuf compiler plugin for generating a JSON dump file for the `protoc` plugin generation request. It is primarily useful for running tests without invoking `protoc` directly (e.g. via the Go testing suite).

## [protoc-gen-dump](https://github.com/sourcegraph/prototools/blob/master/README.dump.md)

protoc-gen-dump is just like `protoc-gen-json` except it dumps output in protobuf format itself. It's much better than `protoc-gen-json` if you can make use of it because there are certain fields (e.g. [extensions](http://godoc.org/github.com/golang/protobuf/protoc-gen-go/descriptor#MethodOptions)) which are explicitly marked as non-JSON-encodable. If you wish to retain e.g. grpc-gateway google.api.http annotations, you'll need to make use of `protoc-gen-dump` instead of `protoc-gen-json`.
