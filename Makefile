# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
ENV_LOCAL_TEST=\
	SUBSCRIPTION=31e9f9a0-9fd2-4294-a0a3-0101246d9700 \
	RESOURCE_GROUP_NAME=aMODMTestb \
	RESOURCE_GROUP_LOCATION=eastus \
	ENV=local \
	PORT=8080
	
apiserver-local: apiserver
	./scripts/apiserver-local.sh

apiserver:
	go build -o ./bin/ ./cmd/apiserver

operator:
	go build -o ./bin/ ./cmd/operator

test-all:
	go test ./...

test-integration:
	$(ENV_LOCAL_TEST) \
	go test -tags=integration ./test -v -count=1 

sdk:
	go build ./sdk

.NOTPARALLEL:

.PHONY: apiserver-local apiserver sdk operator test-all
