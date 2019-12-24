default: generate

generate:
	protoc \
		-I ../third_party/googleapis \
		-I . \
	    --go_out=. \
	    *.proto

clean:
	rm *.pb.go
