package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/protobuf/proto"
	go_descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	protoc2openapi_options "github.com/zhyuri/protoc2openapi/proto"
	"path/filepath"
	"strings"
)

func extractBaseOption(input *descriptor.File, api openapi3.Swagger) (openapi3.Swagger, error) {
	options := input.GetOptions()
	if options == nil {
		return openapi3.Swagger{}, nil
	}
	if !proto.HasExtension(options, protoc2openapi_options.E_Openapi) {
		return openapi3.Swagger{}, nil
	}
	ext, err := proto.GetExtension(options, protoc2openapi_options.E_Openapi)
	if err != nil {
		return openapi3.Swagger{}, err
	}
	if pb, ok := ext.(*protoc2openapi_options.OpenAPI); ok {
		return convertPbToSwagger(pb, api)
	}
	return openapi3.Swagger{}, fmt.Errorf("expect OpenAPI options, got %T", ext)
}

func extractOperation(mthPb *go_descriptor.MethodDescriptorProto) (openapi3.Operation, error) {
	options := mthPb.GetOptions()
	if options == nil {
		return openapi3.Operation{}, nil
	}
	if !proto.HasExtension(options, protoc2openapi_options.E_Operation) {
		return openapi3.Operation{}, nil
	}
	ext, err := proto.GetExtension(options, protoc2openapi_options.E_Operation)
	if err != nil {
		return openapi3.Operation{}, err
	}
	if pb, ok := ext.(*protoc2openapi_options.Operation); ok {
		return convertPbToOperation(pb)
	}
	return openapi3.Operation{}, fmt.Errorf("expect OpenAPI Operation, got %T", ext)
}

// encodeSwagger converts swagger file obj to plugin.CodeGeneratorResponse_File
func encodeSwagger(fileName string, api openapi3.Swagger) (*plugin.CodeGeneratorResponse_File, error) {
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")
	if err := enc.Encode(api); err != nil {
		return nil, err
	}
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	output := fmt.Sprintf("%s.swagger.json", base)
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(output),
		Content: proto.String(formatted.String()),
	}, nil
}

func convertPbToSwagger(pb *protoc2openapi_options.OpenAPI, api openapi3.Swagger) (openapi3.Swagger, error) {
	pbJson, err := json.Marshal(pb)
	if err != nil {
		return openapi3.Swagger{}, err
	}
	err = json.Unmarshal(pbJson, &api)
	return api, err
}

func convertPbToOperation(pb *protoc2openapi_options.Operation) (openapi3.Operation, error) {
	pbJson, err := json.Marshal(pb)
	if err != nil {
		return openapi3.Operation{}, err
	}
	opt := openapi3.Operation{}
	err = json.Unmarshal(pbJson, &opt)
	return opt, err
}
