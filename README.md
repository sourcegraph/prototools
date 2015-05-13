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

The basic syntax is to invoke protoc with `--doc_out`:

```
protoc --doc_out="<OUT_DIR>" input.proto
protoc --doc_out="<OPTIONS>:<OUT_DIR>" input.proto
```

Where `<OPTIONS>` is a comma-seperated list of `key=value` options listed below, for example:

```
protoc --doc_out="template=templates/tmpl.html:doc/" file.proto
```

Would produce documentation using `templates/tmpl.html` inside the `doc/` output directory for `file.proto`.

# Options

| Option     | Description                                                        |
|------------|--------------------------------------------------------------------|
| `template` | Input `.html` `text/template` template file to use for generation. |
