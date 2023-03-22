# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
apiserver-local:
	# build something locally

apiserver:
	go build -o ./bin/ ./cmd/apiserver

sdk:
	go build ./sdk

.NOTPARALLEL:

.PHONY: apiserver-local apiserver sdk