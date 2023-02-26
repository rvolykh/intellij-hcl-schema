requirements: export PBURL ?= "https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protoc-3.15.6-osx-x86_64.zip"
requirements:
	@ mkdir -p $(CURDIR)/.requirements/
	@ echo "Install Protocol Buffer Compiler"
	@ curl -L ${PBURL} -o $(CURDIR)/.requirements/protoc.zip
	@ unzip $(CURDIR)/.requirements/protoc.zip -d $(CURDIR)/.requirements
	@ echo "Install Protocol Buffer Compiler Go plugins"
	@ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	@ go build -o $(CURDIR)/.requirements/bin/protoc-gen-go github.com/golang/protobuf/protoc-gen-go
	@ go build -o $(CURDIR)/.requirements/bin/protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc
.PHONY: requirements

generate: export PATH:=$(CURDIR)/.requirements/bin:${PATH}:$(PATH)
generate:
	@ echo "Generate code"
	@ go generate main.go
.PHONY: generate

build:
	@ mkdir -p bin
	@ go build -o bin/intellij-hcl-schema .
.PHONY: build
