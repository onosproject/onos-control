#!/bin/sh

# SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

proto_imports=".:${GOPATH}/src/github.com/google/protobuf/src:${GOPATH}/src"

protoc -I=$proto_imports --go_out=plugins=grpc:. pkg/northbound/proto/*.proto
