# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
generate:
	./scripts/generate-code.sh

ENV_LOCAL_TEST=\
	SUBSCRIPTION=31e9f9a0-9fd2-4294-a0a3-0101246d9700 \
	RESOURCE_GROUP_NAME=aMODMTestb \
	RESOURCE_GROUP_LOCATION=eastus \
	ENV=local \
	PORT=8080 

clean:
	rm -rf bin
	mkdir -p bin

check-credentials:
ifndef CONTAINER_REGISTRY_USERNAME
	$(error Environment variable CONTAINER_REGISTRY_USERNAME is not set)
endif
ifndef CONTAINER_REGISTRY_PASSWORD
	$(error Environment variable CONTAINER_REGISTRY_PASSWORD is not set)
endif
ifndef CONTAINER_REGISTRY_DEFAULT_SERVER
	$(error Environment variable CONTAINER_REGISTRY_DEFAULT_SERVER is not set)
endif

resolve-registry:
ifndef CONTAINER_REGISTRY_NAMESPACE
CONTAINER_REGISTRY_NAMESPACE := ${CONTAINER_REGISTRY_DEFAULT_NAMESPACE}
endif
ifndef CONTAINER_REGISTRY
CONTAINER_REGISTRY := ${CONTAINER_REGISTRY_DEFAULT_SERVER}/${CONTAINER_REGISTRY_NAMESPACE}
endif
	
apiserver-local:
	./scripts/run-local.sh apiserver

apiserver:
	go build -o ./bin/ ./cmd/apiserver

operator:
	go build -o ./bin/ ./cmd/operator

apiserver-local:
	./scripts/run-local.sh apiserver

operator-local:
	./scripts/run-local.sh operator

test-all:
	go test ./...

test-integration:
	$(ENV_LOCAL_TEST) \
	go test -tags=integration ./test -v -count=1 

sdk:
	go build ./sdk

tools:
	./scripts/build-tools.sh

assemble: apiserver operator 
	./scripts/assemble.sh

.NOTPARALLEL:

.PHONY: apiserver-local apiserver sdk operator test-all generate tools
