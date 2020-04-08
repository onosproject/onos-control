export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_CONTROL_VERSION := latest
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-control ./cmd/onos
	go build -o build/_output/onos ./cmd/onos

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	go test github.com/onosproject/onos-control/pkg/...
	go test github.com/onosproject/onos-control/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-control \
		-w /go/src/github.com/onosproject/onos-control \
		--entrypoint pkg/northbound/proto/compile-protos.sh \
		onosproject/protoc-go:stable

onos-control-base-docker: # @HELP build onos-control base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-control-base:${ONOS_CONTROL_VERSION}
	@rm -rf vendor

onos-control-docker: onos-control-base-docker # @HELP build onos-control Docker image
	docker build . -f build/onos-control/Dockerfile \
		--build-arg ONOS_CONTROL_BASE_VERSION=${ONOS_CONTROL_VERSION} \
		-t onosproject/onos-control:${ONOS_CONTROL_VERSION}

onos-cli-docker: onos-control-base-docker # @HELP build onos-cli Docker image
	docker build . -f build/onos-cli/Dockerfile \
		--build-arg ONOS_CONTROL_BASE_VERSION=${ONOS_CONTROL_VERSION} \
		-t onosproject/onos-cli:${ONOS_CONTROL_VERSION}

onos-control-it-docker: onos-control-base-docker # @HELP build onos-control-integration-tests Docker image
	docker build . -f build/onos-it/Dockerfile \
		--build-arg ONOS_CONTROL_BASE_VERSION=${ONOS_CONTROL_VERSION} \
		-t onosproject/onos-control-integration-tests:${ONOS_CONTROL_VERSION}

# integration: @HELP build and run integration tests
integration: kind
	onit create cluster
	onit add simulator
	onit add simulator
	onit run suite integration-tests


images: # @HELP build all Docker images
images: build onos-control-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-config:${ONOS_CONTROL_VERSION}

all: build images


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-control/onos-control ./cmd/dummy/dummy

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
