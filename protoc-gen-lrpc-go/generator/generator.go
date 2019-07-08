package generator

import (
	"bytes"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strings"
)

func New() *Generator {
	g := new(Generator)
	g.Buffer = new(bytes.Buffer)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)
	return g
}

type Generator struct {
	*bytes.Buffer

	Request  *plugin.CodeGeneratorRequest  // The input.
	Response *plugin.CodeGeneratorResponse // The output.

	Param             map[string]string // Command-line parameters.
	PackageImportPath string            // Go import path of the package we're generating code for
	ImportPrefix      string            // String to prefix to imported package file names.
	ImportMap         map[string]string // Mapping from .proto file name to import path

	Pkg map[string]string // The names under which we import support packages
}

// Error reports a problem, including an error, and exits the program.
func (g *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("[protoc-gen-lrpc-go] error:", s)
	os.Exit(1)
}

// Fail reports a problem and exits the program.
func (g *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("[protoc-gen-lrpc-go] error:", s)
	os.Exit(1)
}

func (g *Generator) Log(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Println("[protoc-gen-lrpc-go] ", s)
}

func (g *Generator) LogObject(v interface{}) {
	s, _ := jsoniter.MarshalToString(v)
	log.Println("[protoc-gen-lrpc-go] ", s)
}

type Package struct {
	Name          string
	Defines       []*descriptor.DescriptorProto
	Service       []*descriptor.ServiceDescriptorProto
}

func (g *Generator) Generate() {
	all := make(map[string]*Package)

	for _, v := range g.Request.ProtoFile {

		pack, ok := all[v.GetPackage()]
		if !ok {
			pack = &Package{}
			pack.Name = v.GetPackage()
			all[pack.Name ] = pack
		}

		for _, vv := range v.MessageType {
			pack.Defines = append(pack.Defines, vv)
		}

		for _, vv := range v.Service {
			pack.Service = append(pack.Service, vv)
		}
	}

	for _, d := range all {
		g.generateGo(d)
	}
}
