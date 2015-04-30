# prototools

This repository holds various Protobuf/gRPC tools.

# protoc-gen-doc

`cmd/protoc-gen-doc` contains a protoc compiler plugin for based documentation generation of `.proto` files. It operates on standard Go `text/template` files (see the `templates` directory) and can produce HTML documentation.

# Installation

First install Go and Protobuf itself, then install the tools using go get:

```
go get -u sourcegraph.com/sourcegraph/prototools/...
```
