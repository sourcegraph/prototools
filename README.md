# prototools

This repository holds various Protobuf/gRPC tools.

# protoc-gen-doc

`cmd/protoc-gen-doc` contains a protoc compiler plugin for based documentation generation of `.proto` files. It operates on standard Go `text/template` files (see the `templates` directory) and can produce HTML documentation.

# Installation

First install Go and Protobuf itself, then install the tools using go get:

```
go get -u sourcegraph.com/sourcegraph/prototools/...
```

# Usage

The basic syntax is to invoke protoc with `--doc_out` where all documentation-related options are passed separated by colons (`:`):

```
protoc --doc_out="<OPTIONS>:<OUT_DIR>" input.proto
```

For example:

```
protoc --doc_out="template=templates/tmpl.html:doc/" file.proto
```

Would produce documentation using `templates/tmpl.html` inside the `doc/` output directory for `file.proto`.
