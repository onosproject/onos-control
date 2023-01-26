// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"github.com/atomix/go-sdk/pkg/test"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-net-lib/pkg/p4utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoresBasics(t *testing.T) {
	client := test.NewClient()

	ctx := context.TODO()

	info, err := p4utils.LoadP4Info("p4info.txt")
	assert.NoError(t, err)

	stores := NewStoreManager(client)

	fooStore, err := stores.Get(ctx, "foo", info)
	assert.NoError(t, err)
	assert.Len(t, fooStore.(*entityStore).tables, 20)
	assert.Equal(t, topo.ID("foo"), fooStore.ID())

	barStore, err := stores.Get(ctx, "bar", info)
	assert.NoError(t, err)
	assert.Len(t, barStore.(*entityStore).tables, 20)
	assert.Equal(t, topo.ID("bar"), barStore.ID())

	foo2Store, err := stores.Get(ctx, "foo", info)
	assert.NoError(t, err)
	assert.Same(t, fooStore, foo2Store)
	assert.Equal(t, topo.ID("foo"), foo2Store.ID())

	assert.Len(t, stores.GetAll(), 2)

	err = stores.Purge(ctx, "foo")
	assert.NoError(t, err)
	err = stores.Purge(ctx, "bar")
	assert.NoError(t, err)
}
