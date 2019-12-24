package generator

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"log"
)

type generator struct {
	reg *descriptor.Registry
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry) *generator {
	return &generator{reg: reg}
}

func (g generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	openAPIs := make(map[string]openapi3.Swagger)
	for _, target := range targets {
		log.Printf("Processing %s", target.GetName())
		openAPI, err := g.generateInstance(target)
		if err != nil {
			return nil, err
		}
		openAPIs[target.GetName()] = openAPI
	}

	var files []*plugin.CodeGeneratorResponse_File
	for fileName, api := range openAPIs {
		f, err := encodeSwagger(fileName, api)
		if err != nil {
			return nil, fmt.Errorf("failed to encode swagger for %s: %s", fileName, err)
		}
		files = append(files, f)
	}
	return files, nil
}

func (g generator) generateInstance(input *descriptor.File) (api openapi3.Swagger, err error) {
	api = openapi3.Swagger{
		OpenAPI: "3.0.2",
	}
	if api, err = extractBaseOption(input, api); err != nil {
		return openapi3.Swagger{}, err
	}
	if err := g.fillServices(input, &api); err != nil {
		return openapi3.Swagger{}, err
	}

	return api, nil
}

func (g generator) fillServices(input *descriptor.File, api *openapi3.Swagger) error {
	components := openapi3.Components{}
	paths := make(openapi3.Paths)
	for _, service := range input.Services {
		for _, method := range service.Methods {
			operation, err := extractOperation(method.MethodDescriptorProto)
			if err != nil {
				return err
			}

			//log.Printf("Method FQMN %s", method.FQMN())
			for i, binding := range method.Bindings {
				log.Printf("[%d] Binding %#v", i, binding)
				log.Printf("[%d] HTTP Method %v", i, binding.HTTPMethod)

				pathItem := &openapi3.PathItem{}
				pathItem.SetOperation(binding.HTTPMethod, &operation)
				// TODO: Fill in actual path
				paths[binding.PathTmpl.Template] = pathItem
			}
		}
	}
	api.Components = components
	api.Paths = paths
	return nil
}
