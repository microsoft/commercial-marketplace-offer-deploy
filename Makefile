# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
generate:
	./scripts/generate-code.sh

apiserver-local: apiserver
	./scripts/apiserver-local.sh

apiserver:
	go build -o ./bin/ ./cmd/apiserver

operator:
	go build -o ./bin/ ./cmd/operator

test-all:
	go test ./...

sdk:
	go build ./sdk

.NOTPARALLEL:

.PHONY: apiserver-local apiserver sdk operator test-all generate
