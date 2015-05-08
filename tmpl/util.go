package tmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// parseProtoPath parses the PROTO_PATH environment variable.
func parseProtoPath() (paths []string) {
	for _, s := range strings.Split(os.Getenv("PROTO_PATH"), " ") {
		paths = append(paths, strings.TrimPrefix(s, "--proto_path="))
	}
	return
}

// findExt iterates through the templates and finds the first file extension it
// can, or returns an empty string if none is found. It should be invoked
// initially with g.Template
func findExt(t *template.Template) string {
	// First this template itself.
	if ext := filepath.Ext(t.Name()); len(ext) > 0 {
		return ext
	}
	// And then the associated templates, recursively.
	for _, tmpl := range t.Templates() {
		if ext := findExt(tmpl); len(ext) > 0 {
			return ext
		}
	}
	return ""
}

// unixPath takes a path, cleans it, and replaces any windows separators (\\)
// with unix ones (/). This is needed because plugin.CodeGeneratorResponse_File
// is defined as having always unix path separators for the file name.
func unixPath(s string) string {
	s = filepath.Clean(s)
	s = strings.Replace(s, "\\", "/", -1)

	// Duplicate clean for trailing slashes that were previously windows ones.
	return filepath.Clean(s)
}

// stripExt strips the extension off the path and returns it.
func stripExt(s string) string {
	ext := filepath.Ext(s)
	if len(ext) > 0 {
		return s[:len(s)-len(ext)]
	}
	return s
}

var rePkgStmt = regexp.MustCompile(`(?m)^\s*package\s+(\w+)\s*;`)

// pkgStmt parses the package statement in a proto file, returning the package
// name, or an empty string if parsing fails.
//
// r must satisfy the io.RuneReader, io.Reader, and io.Seeker interfaces.
func pkgStmt(r interface{}) (pkg string, err error) {
	m := rePkgStmt.FindReaderSubmatchIndex(r.(io.RuneReader))
	if len(m) < 4 {
		return "", nil
	}
	nameStart := m[2]
	nameEnd := m[3]

	// Seek to the start of the package name.
	_, err = r.(io.Seeker).Seek(int64(nameStart), 0)
	if err != nil {
		return "", err
	}

	// Read the package name.
	pkgName := make([]byte, nameEnd-nameStart)
	_, err = io.ReadFull(r.(io.Reader), pkgName)
	if err != nil {
		return "", err
	}
	return string(pkgName), nil
}

var Preload = (&tmplFuncs{}).funcMap()

// cacheItem is a single cache item with a value and a location -- effectively
// it is just used for searching.
type cacheItem struct {
	V interface{}
	L *descriptor.SourceCodeInfo_Location
}

// newTmplFuncs returns an initialized template function structure.
func newTmplFuncs(fd *descriptor.FileDescriptorProto, ext string) *tmplFuncs {
	f := &tmplFuncs{
		f:   fd,
		ext: ext,
	}
	f.pkgNameCache = make(map[string]string, 16)
	f.protoPath = parseProtoPath()
	return f
}

// Functions exposed to templates. The user of the package must first preload
// the FuncMap above for these to be called properly (as they are actually
// closures with context).
type tmplFuncs struct {
	f   *descriptor.FileDescriptorProto
	ext string

	protoPath    []string
	locCache     []cacheItem
	pkgNameCache map[string]string // Map of file names to their package names.
}

// funcMap returns the function map for feeding into templates.
func (f *tmplFuncs) funcMap() template.FuncMap {
	return map[string]interface{}{
		"cleanLabel":     f.cleanLabel,
		"cleanType":      f.cleanType,
		"fieldType":      f.fieldType,
		"urlToType":      f.urlToType,
		"fullyQualified": f.fullyQualified,
		"location":       f.location,
	}
}

// cleanLabel returns the clean (i.e. human-readable / protobuf-style) version
// of a label.
func (f *tmplFuncs) cleanLabel(l *descriptor.FieldDescriptorProto_Label) string {
	switch int32(*l) {
	case 1:
		return "optional"
	case 2:
		return "required"
	case 3:
		return "repeated"
	default:
		panic("unknown label")
	}
}

// cleanType returns the last part of a types name, i.e. for a fully-qualified
// type ".foo.bar.baz" it would return just "baz".
func (f *tmplFuncs) cleanType(path string) string {
	split := strings.Split(path, ".")
	return split[len(split)-1]
}

// fieldType returns the clean (i.e. human-readable / protobuf-style) version
// of a field type.
func (f *tmplFuncs) fieldType(field *descriptor.FieldDescriptorProto) string {
	if field.TypeName != nil {
		return f.cleanType(*field.TypeName)
	}
	switch int32(*field.Type) {
	case 1:
		return "double"
	case 2:
		return "float"
	case 3:
		return "int64"
	case 4:
		return "uint64"
	case 5:
		return "int32"
	case 6:
		return "fixed64"
	case 7:
		return "fixed32"
	case 8:
		return "bool"
	case 9:
		return "string"
	case 10:
		return "group"
	case 11:
		return "message"
	case 12:
		return "bytes"
	case 13:
		return "uint32"
	case 14:
		return "enum"
	case 15:
		return "sfixed32"
	case 16:
		return "sfixed64"
	case 17:
		return "sint32"
	case 18:
		return "sint64"
	default:
		panic("unknown type")
	}
}

// urlToType returns a URL to the documentation file for the given type. The
// input type path can be either fully-qualified or not, regardless, the URL
// returned will always have a fully-qualified hash.
func (f *tmplFuncs) urlToType(typePath string) (string, error) {
	typePath = f.fullyQualified(typePath)

	// Resolve the package path for the type.
	pkg := strings.Split(typePath, ".")[1]
	pkgPath, err := f.resolvePkgPath(pkg)
	if err != nil {
		return "", err
	}
	if pkgPath == "" {
		return "", nil
	}

	// Make the path relative to this documentation files directory and then swap
	// the extension out.
	basePath := filepath.Dir(*f.f.Name)
	rel, _ := filepath.Rel(basePath, pkgPath)
	rel = stripExt(rel) + f.ext
	return fmt.Sprintf("%s#%s", rel, typePath), nil
}

// fullyQualified returns the fully qualified path for the given type path.
//
// TODO(slimsag): this is incomplete as it assumes the scope is only relative
// to the package, i.e. for ".foo.bar.baz" fullyQualified("baz") would return
// ".foo.baz" incorrectly. Handling such cases requires more extensive C++
// style scope crawling.
func (f *tmplFuncs) fullyQualified(typePath string) string {
	if typePath[0] == '.' {
		return typePath
	}

	// Not fully-qualified.
	pkg := stripExt(filepath.Base(*f.f.Name))
	return fmt.Sprintf(".%s.%s", pkg, typePath)
}

// pkgName parses the named file path for a package statement and returns the
// package name.
func (f *tmplFuncs) pkgName(path string) (string, error) {
	// Grab the name from cache if we already have it.
	cached, ok := f.pkgNameCache[path]
	if ok {
		return cached, nil
	}

	// Parse for the package name and put into cache.
	var (
		file *os.File
		err  error
	)
	for _, includeDir := range f.protoPath {
		file, err = os.Open(filepath.Join(includeDir, path))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		defer file.Close()
		break
	}
	if file == nil {
		return "", fmt.Errorf("not found in PROTO_PATH: %s", path)
	}

	// *os.File is -- but bufio.Reader isn't -- an io.Seeker, so we end up reading
	// the whole file into memory anyway.
	//
	// TODO(slimsag): adopt a line-by-line approach instead.
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	pkg, err := pkgStmt(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if pkg == "" {
		// It's the file name then,
		pkg = stripExt(filepath.Base(path))
	}
	f.pkgNameCache[path] = pkg
	return pkg, nil
}

// resolvePkgPath resolves the named protobuf package, returning it's file
// path.
func (f *tmplFuncs) resolvePkgPath(pkg string) (string, error) {
	// Test this proto file itself:
	name, err := f.pkgName(*f.f.Name)
	if err != nil {
		return "", err
	}
	if name == pkg {
		return *f.f.Name, nil
	}

	// Test each dependency:
	for _, p := range f.f.Dependency {
		name, err = f.pkgName(p)
		if err != nil {
			return "", err
		}
		if name == pkg {
			return p, nil
		}
	}
	return "", nil
}

// location returns the source code info location for the generic AST-like node
// from the descriptor package.
func (f *tmplFuncs) location(x interface{}) *descriptor.SourceCodeInfo_Location {
	// Validate that we got a sane type from the template.
	pkgPath := reflect.Indirect(reflect.ValueOf(x)).Type().PkgPath()
	if pkgPath != "" && pkgPath != "github.com/golang/protobuf/protoc-gen-go/descriptor" {
		panic("expected descriptor type; got " + fmt.Sprintf("%q", pkgPath))
	}

	// If the location cache is empty; we build it now.
	if f.locCache == nil {
		for _, loc := range f.f.SourceCodeInfo.Location {
			f.locCache = append(f.locCache, cacheItem{
				V: f.walkPath(loc.Path),
				L: loc,
			})
		}
	}
	return f.findCachedItem(x)
}

// findCachedItem finds and returns a cached location for x.
func (f *tmplFuncs) findCachedItem(x interface{}) *descriptor.SourceCodeInfo_Location {
	for _, i := range f.locCache {
		if i.V == x {
			return i.L
		}
	}
	return nil
}

// walkPath walks through the root node (the f.f file) descending down the path
// until it is resolved, at which point the value is returned.
func (f *tmplFuncs) walkPath(path []int32) interface{} {
	if len(path) == 0 {
		return f.f
	}
	var (
		walker func(id int, v interface{}) bool
		found  interface{}
		target = int(path[0])
	)
	path = path[1:]
	walker = func(id int, v interface{}) bool {
		if id != target {
			return true
		}
		if len(path) == 0 {
			found = v
			return false
		}
		target = int(path[0])
		path = path[1:]
		f.protoFields(reflect.ValueOf(v), walker)
		return false
	}
	f.protoFields(reflect.ValueOf(f.f), walker)
	return found
}

// protoFields invokes fn with the protobuf tag ID and its in-memory Go value
// given a descriptor node type. It stops invoking fn when it returns false.
func (f *tmplFuncs) protoFields(node reflect.Value, fn func(id int, v interface{}) bool) {
	indirect := reflect.Indirect(node)

	switch indirect.Kind() {
	case reflect.Slice:
		for i := 0; i < indirect.Len(); i++ {
			if !fn(i, indirect.Index(i).Interface()) {
				return
			}
		}

	case reflect.Struct:
		// Iterate each field.
		for i := 0; i < indirect.NumField(); i++ {
			// Parse the protobuf tag for the ID, e.g. the 49 in:
			// "bytes,49,opt,name=foo,def=hello!"
			tag := indirect.Type().Field(i).Tag.Get("protobuf")
			fields := strings.Split(tag, ",")
			if len(fields) < 2 {
				continue // too few fields
			}

			// Parse the tag ID.
			tagID, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			if !fn(tagID, indirect.Field(i).Interface()) {
				return
			}
		}
	}
}
