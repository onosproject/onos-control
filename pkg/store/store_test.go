// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"github.com/atomix/go-sdk/pkg/test"
	"testing"
)

func TestStoreBasics(t *testing.T) {
	_ = test.NewClient()
}

/*
	//ctx := context.TODO()

	info, err := p4utils.LoadP4Info("p4info.txt")
	assert.NoError(t, err)
	t.Log(info)
	//store, err := NewEntityStore(ctx, client, topo.ID("id"), info)
	//assert.NoError(t, err)
	//
	//es := store.(*entityStore)
	//assert.Len(t, es.tables, 5)

*/
