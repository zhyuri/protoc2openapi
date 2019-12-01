package generator

import (
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

type generator struct {
	reg *descriptor.Registry
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry) *generator {
	return &generator{reg: reg}
}

func (g generator) Generate(targets []*descriptor.File) ([]*plugin_go.CodeGeneratorResponse_File, error) {
	// TODO: implement
	return nil, nil
}
