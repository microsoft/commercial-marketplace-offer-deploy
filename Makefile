# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
apiserver-local: apiserver
	./scripts/apiserver-local.sh

apiserver:
	go build -o ./bin/ ./cmd/apiserver

apiserver-test: 
	go test ./cmd/apiserver...

operator:
	go build -o ./bin/ ./cmd/operator

operator-test:
	go test ./cmd/operator...

sdk:
	go build ./sdk

.NOTPARALLEL:

.PHONY: apiserver-local apiserver sdk
