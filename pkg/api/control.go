// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package api contains definitions of the library interfaces made available to the applications
package api

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/topo"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

// DeviceControl is an abstraction of an entity allowing control over the
// forwarding behavior of a single device
type DeviceControl interface {
	// Read receives a query and returns back all requested control entries on the given channel
	Read(ctx context.Context, entities *[]p4api.Entity, ch chan<- []*p4api.Entity) error

	// Write applies a set of updates to the device
	Write(ctx context.Context, request *[]p4api.Update) error

	// EmitPacket requests emission of the specified packet onto the data-plane
	EmitPacket(ctx context.Context, packetOut *p4api.PacketOut) error

	// HandlePackets starts handling the packet-in message using the supplied channel and packet handler
	HandlePackets(ch chan<- *p4api.PacketIn, handler *PacketHandler)

	// Pipeline returns the P4 information describing the high-level device pipeline
	Pipeline() *p4info.P4Info

	// Version returns the P4Runtime version of the target
	Version() string
}

// PacketHandler is an abstraction of an entity capable of handling an incoming packet-in
type PacketHandler interface {
	// Handle handles the given packet-in messages
	Handle(packetIn *p4api.PacketIn) error
}

// Devices is an abstraction of an entity capable of tracking device control contexts
// of multiple devices on behalf of the control application.
type Devices interface {
	// Add requests creation of a new device flow control context using its P4Runtime connection endpoint
	Add(ctx context.Context, id topo.ID, p4rtEndpoint string, translator PipelineTranslator) (DeviceControl, error)

	// Remove requests removal of device control context
	Remove(id topo.ID)

	// Get the device flow control entity by its ID
	Get(id topo.ID) DeviceControl

	// GetAll returns all device flow control entities presently registered with the manager
	GetAll() []DeviceControl
}
