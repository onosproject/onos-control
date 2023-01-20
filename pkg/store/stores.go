// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package store contains implementation of the store for persisting pipeline entities
package store

import (
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/onosproject/onos-api/go/onos/topo"
)

// Stores is an abstraction of an entity capable of tracking multiple entity stores created
// on behalf of different devices.
type Stores interface {
	// Get returns an entity store for the specified device. A new one will be created if it doesn't
	// already exist. Returns error if the backing storage system is unavailable.
	Get(id topo.ID) (EntityStore, error)

	// GetAll returns the list of all currently registered entity stores. It returns only those store
	// that have been retrieved and/or implicitly created via retrieval during this run-time, meaning
	// that this list itself is not persisted.
	GetAll() []EntityStore

	// Purge purges the given device entity store and all the data within it.
	Purge(id topo.ID) error
}

type storeManager struct {
	Stores
	stores map[topo.ID]EntityStore
}

// NewStoreManager creates a new stores manager
func NewStoreManager(client primitive.Client) Stores {
	// TODO: Implement me
	return &storeManager{stores: make(map[topo.ID]EntityStore, 0)}
}

// Get returns an entity store for the specified device.
func (sm *storeManager) Get(id topo.ID) (EntityStore, error) {
	// TODO: Implement me
	return nil, nil
}

// GetAll returns the list of all currently registered entity stores.
func (sm *storeManager) GetAll() []EntityStore {
	stores := make([]EntityStore, 0, len(sm.stores))
	for _, store := range sm.stores {
		stores = append(stores, store)
	}
	return stores
}

// Purge purges the given device entity store and all the data within it.
func Purge(id topo.ID) error {
	// TODO: Implement me
	return nil
}
