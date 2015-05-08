# prototools

This repository holds various Protobuf/gRPC tools.

# protoc-gen-doc

`cmd/protoc-gen-doc` contains a protoc compiler plugin for based documentation generation of `.proto` files. It operates on standard Go `html/template` files (see the `templates` directory) and can produce HTML documentation.

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

# proto_path issues

You will quickly find that the generator fails if you run it with a `--proto_path` argument. This is because `protoc` does not pass it's command-line arguments onto generator programs (and `protoc-gen-doc` needs to know the `--proto_path` arguments in order to resolve symbols appropriately). To workaround this issue you must export the `$PROTO_PATH` environment variable, thus you can run `protoc` as just:

```
PROTO_PATH="--proto_path=/my/custom/path" protoc $(PROTO_ARGS) --doc_out="template=templates/tmpl.html:doc/" file.proto
```
