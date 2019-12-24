
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

default: build

generate:
	(cd ./proto && make generate)

.PHONY: build
build: generate
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) main.go

.PHONY: test
test: build
	mkdir -p $(OUTPUT_DIR)/openapi/
	@PATH=$(PATH):$(BIN_DIR) \
	$(PROTOC) \
		$(PROTO_INCLUDES) \
	    --openapi_out=$(OUTPUT_DIR)/openapi/ \
	    example/*.proto

clean:
	rm -rf build/*
	(cd ./proto && make clean)
