// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package controller implements the core reconciliation controller tying together API, stores, translators and SB
package controller

import (
	"context"
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-control/pkg/api"
	"github.com/onosproject/onos-control/pkg/store"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
	"sync"
)

type devicesController struct {
	api.Devices
	role   *p4api.Role
	stores store.Stores

	mu      sync.RWMutex
	devices map[topo.ID]*deviceController
}

// NewController creates a new controller for device control contexts using the supplied role descriptor
// and pipeline translator
func NewController(role *p4api.Role, client primitive.Client) api.Devices {
	return &devicesController{
		role:   role,
		stores: store.NewStoreManager(client),
	}
}

// Add requests creation of a new device flow control context using its P4Runtime connection endpoint
func (c *devicesController) Add(ctx context.Context, id topo.ID, p4rtEndpoint string, translator api.PipelineTranslator) (api.DeviceControl, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if d, ok := c.devices[id]; ok {
		return d, nil
	}

	s, err := c.stores.Get(ctx, id, translator.FromPipeline())
	if err != nil {
		return nil, err
	}
	return newDeviceController(id, p4rtEndpoint, s, translator), nil
}

// Remove requests removal of device control context
func (c *devicesController) Remove(id topo.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.devices, id)
}

// Get the device flow control entity by its ID
func (c *devicesController) Get(id topo.ID) api.DeviceControl {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.devices[id]
}

// GetAll returns all device flow control entities presently registered with the manager
func (c *devicesController) GetAll() []api.DeviceControl {
	c.mu.RLock()
	defer c.mu.RUnlock()
	devices := make([]api.DeviceControl, 0, len(c.devices))
	for _, d := range c.devices {
		devices = append(devices, d)
	}
	return devices
}
