// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package store contains implementation of the store for persisting pipeline entities
package store

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/topo"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

// EntityStore is an abstraction of a store capable of maintaining various P4 entities
type EntityStore interface {
	// ID returns the ID of the device whose control entities it persists.
	ID() topo.ID

	// Read accepts a query in form of a list of partially populated entities and returns any
	// matching entities on the specified channel
	Read(ctx context.Context, query []*p4api.Entity, ch chan<- p4api.Entity) error

	// Write persists the specified list of updates.
	Write(ctx context.Context, updates []*p4api.Update) error
}

type entityStore struct {
	EntityStore
	id topo.ID
	// TODO: Insert Atomix primitives to track table, group, meter, etc. entries
}

// NewEntityStore creates a new P4 entity store for the specified device
func NewEntityStore(id topo.ID) (EntityStore, error) {
	// TODO: Implement me
	return &entityStore{id: id}, nil
}

// ID returns the ID of the device whose control entities it persists.
func (s *entityStore) ID() topo.ID {
	return s.id
}

// Read accepts a query in form of a list of partially populated entities and returns any
// matching entities on the specified channel
func (s *entityStore) Read(ctx context.Context, query []*p4api.Entity, ch chan<- p4api.Entity) error {
	// TODO: Implement me
	return nil
}

// Write persists the specified list of updates.
func (s *entityStore) Write(ctx context.Context, updates []*p4api.Update) error {
	// TODO: Implement me
	return nil
}
