
OUTPUT_DIR := build
BIN_DIR := $(OUTPUT_DIR)/bin
# Name must match protoc-gen-{NAME}
# https://developers.google.com/protocol-buffers/docs/reference/cpp/google.protobuf.compiler.plugin
BIN := $(BIN_DIR)/protoc-gen-openapi

PROTOC := protoc
PROTO_INCLUDES := \
	-I third_party/googleapis \
	-I proto \
	-I .
PROTO_OUTPUT := ${OUTPUT_DIR}/proto-gen

default: build

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) main.go

.PHONY: test
test: build
	mkdir -p $(OUTPUT_DIR)/openapi/
	$(PROTOC) \
		$(PROTO_INCLUDES) \
	    --openapi_out=$(OUTPUT_DIR)/openapi/ \
	    example/*.proto

clean:
	rm -rf build/*
