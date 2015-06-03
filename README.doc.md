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

| Option         | Default                     | Description                                                        |
|----------------|-----------------------------|--------------------------------------------------------------------|
| `template`     | `templates/tmpl.html`       | Input `.html` `html/template` template file to use for generation. |
| `root`         | (current working directory) | Root directory path to prefix all generated URLs with.             |
| `filemap`      | none                        | A XML filemap, which specifies how output files are generated.     |
| `dump-filemap` | none                        | Dump the executed filemap template to the given filepath.          |
| `apihost`      | none                        | (grpc-gateway) API host base URL (e.g. `http://api.mysite.com/`)   |

The `template` and `filemap` options are exclusive (only one may be used at a time).

## Templates

Template files are standard Go `html/template` files, and as such their documentation can be found [in that package](https://golang.org/pkg/html/template).

## File Maps

In many cases producing a single output `.html` file for a single input `.proto` file is not desired, often producing very verbose or long web pages. Because protoc-gen-doc doesn't really know how you want your documentation laid out on the file-system (and does not want to restrict you), we offer templated XML file maps.

Say for example that we wanted to produce documentation for three gRPC services, all of which are declared inside of a `organization/services.proto` file:

- "Producer" -> `out_dir/organization/service-producer.html`
- "Consumer" -> `out_dir/organization/service-consumer.html`
- "Trader"   -> `out_dir/organization/service-trader.html`

And also an `out_dir/index.html` file which will link to the three services.

In order to produce the above three (plus index) files for the single `services.proto` file, we will use an XML filemap that follows the syntax of:

```
<FileMap>
    <Generate>
        <Template>path/to/template.html</Template>
        <Target>path/to/target.proto</Target> <!-- optional -->
        <Output>path/to/output.html</Output>
        <Includes>
            <Include>a.tmpl</Include>
            <Include>b.tmpl</Include>
        </Includes>
        <Data>
            <Item><Key>myKey1</Key><Value>myValue1</Value></Item>
            <Item><Key>myKey2</Key><Value>myValue2</Value></Item>
        </Data>
    </Generate>
    <Generate>
        ...
    </Generate>
</FileMap>
```

Where `<Template>` and `<Include>` paths are relative to the XML filemap directory, `<Output>` paths are relative to the output directory, and `<Target>` paths are relative to `--proto_path` directories.

Thus we would write:

```
<FileMap>
    <!-- Index Page (notice the lack of target tag) -->
    <Generate>
        <Template>index.html</Template>
        <Output>index.html</Output>
    </Generate>

    <!-- Service Pages -->
    <Generate>
        <Template>service.html</Template>
        <Target>organization/services.proto</Target>
        <Output>organization/service-producer.html</Output>
        <Data>
            <Item><Key>Service</Key><Value>Producer</Value></Item>
        </Data>
    </Generate>
    <Generate>
        <Template>service.html</Template>
        <Target>organization/services.proto</Target>
        <Output>organization/service-consumer.html</Output>
        <Data>
            <Item><Key>Service</Key><Value>Consumer</Value></Item>
        </Data>
    </Generate>
    <Generate>
        <Template>service.html</Template>
        <Target>organization/services.proto</Target>
        <Output>organization/service-trader.html</Output>
        <Data>
            <Item><Key>Service</Key><Value>Trader</Value></Item>
        </Data>
    </Generate>
</FileMap>
```

As you can see, for just three types it quickly becomes very cumbersome. For this reason filemaps _are themselves templates_ (this is also the rational for choosing XML over e.g. JSON), so we can write the above instead using `html/template` syntax:

```
<FileMap>
    <!-- Index Page (notice the lack of target tag) -->
    <Generate>
        <Template>index.html</Template>
        <Output>index.html</Output>
    </Generate>

{{$serviceTemplate := "service.html"}}
{{range $f := .ProtoFile}}
    {{range $s := .Service}}
        <Generate>
            <Template>{{$serviceTemplate}}</Template>
            <Target>{{$f.Name}}</Target>
            <Output>{{dir $f.Name}}/{{$s.Name}}{{ext $serviceTemplate}}</Output>
            <Data>
                <Item><Key>Service</Key><Value>{{$s.Name}}</Value></Item>
            </Data>
        </Generate>
    {{end}}
{{end}}
</FileMap>
```

Which is to say: for every service type (`range $s := .Service`) in every protobuf input file (`range $f := .ProtoFile` and `<Target>{{$f.Name}}</Target>`) generate a output file using the Go `html/template` (`{{$serviceTemplate}}`) placing output in the directory the protobuf file is in (`{{$f.Name}}`, `organization` in our example), with the service types name (`{{$s.Name}}`, or `Producer` `Consumer` `Trader` above) with the extension of the `$serviceTemplate` file (`.html`), and when each `<Template>` is executed pass along a map with the given keys/values (the service template can then selectively render _just that service type_).

For debugging purposes, you can use the `dump-filemap` option which will execute the template and dump the resulting XML out to a file.

## Issues

If you run into trouble or have questions, please [open an issue](https://github.com/sourcegraph/prototools/issues/new).
