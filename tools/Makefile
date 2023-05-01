
clean:
	rm -rf ../bin/testharness

build:
	../scripts/build-tools.sh

run:
	../scripts/run-local.sh testharness $(build)

test:
	go test ./...

.NOTPARALLEL:

.PHONY: clean build run test
