# SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0


export CGO_ENABLED=1
export GO111MODULE=on

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then mkdir -p build && cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

.PHONY: build

build: # @HELP build the Go binaries (default)
build:
	go build github.com/onosproject/onos-control/pkg/...

mod-update: # @HELP Download the dependencies to the vendor folder
	go mod tidy
	go mod vendor
mod-lint: mod-update # @HELP ensure that the required dependencies are in place
	# dependencies are vendored, but not committed, go.sum is the only thing we need to check
	bash -c "diff -u <(echo -n) <(git diff go.sum)"


test: # @HELP run the unit tests and source code validation  producing a golang style report
test: mod-lint build linters license
	go test -race github.com/onosproject/onos-control/pkg/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: mod-lint build linters license jenkins-tools
	TEST_PACKAGES=github.com/onosproject/onos-control/pkg/... ./build/build-tools/build/jenkins/make-unit

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION}

jenkins-publish: jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/build-tools/release-merge-commit

all: test

clean:: # @HELP remove all the build artifacts
	go clean -testcache github.com/onosproject/onos-control/...
