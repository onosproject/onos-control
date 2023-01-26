// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package controller implements the core reconciliation controller tying together API, stores, translators and SB
package controller

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-control/pkg/api"
	"github.com/onosproject/onos-control/pkg/store"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

type deviceController struct {
	api.DeviceControl
	id         topo.ID
	endpoint   string
	translator api.PipelineTranslator
	version    string
}

func newDeviceController(id topo.ID, endpoint string, store store.EntityStore, translator api.PipelineTranslator) api.DeviceControl {
	return &deviceController{
		id:         id,
		endpoint:   endpoint,
		translator: translator,
	}
}

// Read receives a query and returns back all requested control entries on the given channel
func (d *deviceController) Read(ctx context.Context, entities *[]p4api.Entity, ch chan<- []*p4api.Entity) error {
	// TODO: Implement me
	return nil
}

// Write applies a set of updates to the device
func (d *deviceController) Write(ctx context.Context, request *[]p4api.Update) error {
	// TODO: Implement me
	// Write to the store; notify reconciler to translate and contact southbound
	return nil
}

// EmitPacket requests emission of the specified packet onto the data-plane
func (d *deviceController) EmitPacket(ctx context.Context, packetOut *p4api.PacketOut) error {
	// TODO: Implement me
	// Go directly to the southbound
	return nil
}

// HandlePackets starts handling the packet-in message using the supplied channel and packet handler
func (d *deviceController) HandlePackets(ch chan<- *p4api.PacketIn, handler *api.PacketHandler) {
	// Pass the channel and handler directly to the southbound
}

// Pipeline returns the P4 information describing the high-level device pipeline
func (d *deviceController) Pipeline() *p4info.P4Info {
	return d.translator.FromPipeline()
}

// Version returns the P4Runtime version of the target
func (d *deviceController) Version() string {
	if d.version == "" {
		// TODO: Implement me
		// Issue capabilities request using southbound and capture version here
		d.version = "unknown"
	}
	return d.version
}
