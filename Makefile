# SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_CONTROL_VERSION := latest
ONOS_BUILD_VERSION := stable

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-control ./cmd/onos
	go build -o build/_output/onos ./cmd/onos

test: # @HELP run the unit tests and source code validation
test: build deps linters license
	go test github.com/onosproject/onos-control/pkg/...
	go test github.com/onosproject/onos-control/cmd/...

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

images: # @HELP build all Docker images
images: build onos-control-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-config:${ONOS_CONTROL_VERSION}

all: build images

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-control/onos-control ./cmd/dummy/dummy

