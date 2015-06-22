# protoc-gen-dump

`cmd/protoc-gen-dump` contains a protoc compiler plugin for generating protobuf-encoded dumps of `protoc` plugin generation requests, i.e. the literal request to the `protoc-gen-dump` command from `protoc`.

As it is literally the protobuf-encoded form of [plugin.CodeGeneratorRequest](https://sourcegraph.com/github.com/golang/protobuf@056d5ce64f754d9919f5d66da0735951b4a0e138/.tree/protoc-gen-go/plugin/plugin.pb.go#def=/github.com/golang/protobuf@056d5ce64f754d9919f5d66da0735951b4a0e138/.GoPackage/github.com/golang/protobuf/protoc-gen-go/plugin/.def/CodeGeneratorRequest&startbyte=699&endbyte=2077), it can be unmarshaled into the same exact structure which is useful for writing tests without invoking `protoc` directly (and _attemping_ to detect plugin failure through the compiler proxy).

## Installation

First install Go and Protobuf itself, then install the plugin using go get:

```
go get -u sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-dump
```

## Usage

Simply invoke `protoc` with the `--dump_out` command-line parameter:

```
protoc --dump_out="<OUT_DIR>" input.proto
protoc --dump_out="<OPTIONS>:<OUT_DIR>" input.proto
```

Where `<OPTIONS>` is a comma-seperated list of `key=value,key2=value2` options (which are listed in detail below). For example:

```
protoc --dump_out="dump/" a.proto b.proto
```

Would produce a protobuf-encoded dump file (`out.dump`) for both `a.proto` and `b.proto` inside the `dump/` directory.

## Options

| Option   | Default           | Description                                           |
|----------|-------------------|-------------------------------------------------------|
| `out`    | `out.dump`        | Output file name (_not directory_).                   |

## Issues

If you run into trouble or have questions, please [open an issue](https://github.com/sourcegraph/prototools/issues/new).
