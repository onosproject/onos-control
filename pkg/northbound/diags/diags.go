// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package diags implements the diagnostic gRPC service for the control subsystem.
package diags

import (
	"github.com/onosproject/onos-control/pkg/northbound"
	"github.com/onosproject/onos-control/pkg/northbound/proto"
	"google.golang.org/grpc"
)

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	proto.RegisterControlDiagsServer(r, Server{})
}

// Server implements the gRPC service for diagnostic facilities.
type Server struct {
}
