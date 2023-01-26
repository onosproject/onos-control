// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package store contains implementation of the store for persisting pipeline entities
package store

import (
	"context"
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	"sync"
)

var log = logging.GetLogger("store")

// Stores is an abstraction of an entity capable of tracking multiple entity stores created
// on behalf of different devices.
type Stores interface {
	// Get returns an entity store for the specified device. A new one will be created if it doesn't
	// already exist. Returns error if the backing storage system is unavailable.
	Get(ctx context.Context, id topo.ID, info *p4info.P4Info) (EntityStore, error)

	// GetAll returns the list of all currently registered entity stores. It returns only those store
	// that have been retrieved and/or implicitly created via retrieval during this run-time, meaning
	// that this list itself is not persisted.
	GetAll() []EntityStore

	// Purge purges the given device entity store and all the data within it.
	Purge(ctx context.Context, id topo.ID) error
}

type storeManager struct {
	Stores
	mu     sync.RWMutex
	client primitive.Client
	stores map[topo.ID]EntityStore
}

// NewStoreManager creates a new stores manager
func NewStoreManager(client primitive.Client) Stores {
	return &storeManager{client: client, stores: make(map[topo.ID]EntityStore, 0)}
}

// Get returns an entity store for the specified device.
func (sm *storeManager) Get(ctx context.Context, id topo.ID, info *p4info.P4Info) (EntityStore, error) {
	sm.mu.RLock()
	store, ok := sm.stores[id]
	sm.mu.RUnlock()
	if ok {
		return store, nil
	}

	if info == nil {
		return nil, errors.NewNotFound("Store not found for %s", id)
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	store, ok = sm.stores[id]
	if ok {
		return store, nil
	}

	var err error
	log.Infof("Creating store %s", id)
	store, err = NewEntityStore(ctx, sm.client, id, info)
	if err != nil {
		return nil, errors.FromAtomix(err)
	}
	sm.stores[id] = store
	return store, nil
}

// GetAll returns the list of all currently registered entity stores.
func (sm *storeManager) GetAll() []EntityStore {
	stores := make([]EntityStore, 0, len(sm.stores))
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, store := range sm.stores {
		stores = append(stores, store)
	}
	return stores
}

// Purge purges the given device entity store and all the data within it.
func (sm *storeManager) Purge(ctx context.Context, id topo.ID) error {
	store, err := sm.Get(ctx, id, nil)
	if err != nil {
		return err
	}
	log.Infof("Purging store %s", id)
	return store.(*entityStore).Purge(ctx)
}
