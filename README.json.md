# protoc-gen-json

`cmd/protoc-gen-json` contains a protoc compiler plugin for generating JSON dumps of `protoc` plugin generation requests, i.e. the literal request to the `protoc-gen-json` command from `protoc`.

As it is literally the JSON form of [plugin.CodeGeneratorRequest](https://sourcegraph.com/github.com/golang/protobuf@056d5ce64f754d9919f5d66da0735951b4a0e138/.tree/protoc-gen-go/plugin/plugin.pb.go#def=/github.com/golang/protobuf@056d5ce64f754d9919f5d66da0735951b4a0e138/.GoPackage/github.com/golang/protobuf/protoc-gen-go/plugin/.def/CodeGeneratorRequest&startbyte=699&endbyte=2077), it can be unmarshaled into the same exact structure which is useful for writing tests without invoking `protoc` directly (and _attemping_ to detect plugin failure through the compiler proxy).

## Installation

First install Go and Protobuf itself, then install the plugin using go get:

```
go get -u sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-json
```

## Usage

Simply invoke `protoc` with the `--json_out` command-line parameter:

```
protoc --json_out="<OUT_DIR>" input.proto
protoc --json_out="<OPTIONS>:<OUT_DIR>" input.proto
```

Where `<OPTIONS>` is a comma-seperated list of `key=value,key2=value2` options (which are listed in detail below). For example:

```
protoc --json_out="dump/" a.proto b.proto
```

Would produce a JSON dump file (`out.json`) for both `a.proto` and `b.proto` inside the `dump/` directory.

## Options

| Option   | Default           | Description                                           |
|----------|-------------------|-------------------------------------------------------|
| `out`    | `out.json`        | Output file name (_not directory_).                   |
| `indent` | `  ` (two spaces) | Indention string to use for each nested JSON element. |

## Issues

If you run into trouble or have questions, please [open an issue](https://github.com/sourcegraph/prototools/issues/new).
