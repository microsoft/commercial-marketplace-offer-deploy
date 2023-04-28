# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
generate:
	./scripts/generate-code.sh

clean:
	rm -rf bin
	mkdir -p bin

apiserver:
	go build -o ./bin/ ./cmd/apiserver

operator:
	go build -o ./bin/ ./cmd/operator

apiserver-local:
	./scripts/run-local.sh apiserver

operator-local:
	./scripts/run-local.sh operator

# Builds docker container, starts ngrok in the background, and 
# calls docker compose up with the public NGROK endpoint for MODM to receive event messages from Azure
run-local:
	./scripts/run-local.sh docker $(build)

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

.PHONY: run apiserver-local operator-local apiserver sdk operator test-all generate tools
