# protoc-gen-doc

`cmd/protoc-gen-doc` contains a protoc compiler plugin for documentation generation of `.proto` files. It operates on standard Go `html/template` files (see the `templates` directory) and as such can produce HTML documentation.

## Installation

First install Go and Protobuf itself, then install the plugin using go get:

```
go get -u sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-doc
```

## Usage

Simply invoke `protoc` with the `--doc_out` command-line parameter:

```
protoc --doc_out="<OUT_DIR>" input.proto
protoc --doc_out="<OPTIONS>:<OUT_DIR>" input.proto
```

Where `<OPTIONS>` is a comma-seperated list of `key=value,key2=value2` options (which are listed in detail below). For example:

```
protoc --doc_out="doc/" file.proto
```

Would produce documentation for `file.proto` inside the `doc/` directory using the template `templates/tmpl.html` HTML template file.

## Options

| Option     | Default               | Description                                                        |
|------------|-----------------------|--------------------------------------------------------------------|
| `template` | `templates/tmpl.html` | Input `.html` `text/template` template file to use for generation. |

## Issues

If you run into trouble or have questions, please [open an issue](https://github.com/sourcegraph/prototools/issues/new).
