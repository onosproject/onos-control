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

func TestStoreBasics(t *testing.T) {
	client := test.NewClient()

	ctx := context.TODO()

	info, err := p4utils.LoadP4Info("p4info.txt")
	assert.NoError(t, err)

	store, err := NewEntityStore(ctx, client, topo.ID("id"), info)
	assert.NoError(t, err)

	es := store.(*entityStore)
	assert.Len(t, es.tables, 20)
}
