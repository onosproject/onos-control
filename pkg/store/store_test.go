// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"github.com/atomix/go-sdk/pkg/test"
	"github.com/onosproject/onos-net-lib/pkg/p4utils"
	testutils "github.com/onosproject/onos-net-lib/pkg/test"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestStoreBasics(t *testing.T) {
	client := test.NewClient()

	ctx := context.TODO()

	info, err := p4utils.LoadP4Info("../../test/p4info.txt")
	assert.NoError(t, err)

	store, err := NewEntityStore(ctx, client, "foo", info)
	assert.NoError(t, err)

	es := store.(*entityStore)
	assert.Len(t, es.tables, 20)

	// Generate a slew of random updates and store them
	updates := generateRandomUpdates(info)
	err = store.Write(ctx, updates)
	assert.NoError(t, err)

	// Query all tables
	query := []*p4api.Entity{
		{Entity: &p4api.Entity_TableEntry{TableEntry: &p4api.TableEntry{}}},
	}
	_ = readEntries(ctx, t, store, query, len(updates))

	// Query a specific table
	tableID := info.Tables[0].Preamble.Id
	query = []*p4api.Entity{
		{Entity: &p4api.Entity_TableEntry{TableEntry: &p4api.TableEntry{TableId: tableID}}},
	}
	// Count expected updates for this table
	tupdates := make([]*p4api.Update, 0, len(updates))
	for _, update := range updates {
		if update.GetEntity().GetTableEntry().TableId == tableID {
			tupdates = append(tupdates, update)
		}
	}
	entities := readEntries(ctx, t, store, query, len(tupdates))

	// Remove the first entity
	deletes := []*p4api.Update{{Type: p4api.Update_DELETE, Entity: entities[0]}}
	err = store.Write(ctx, deletes)
	assert.NoError(t, err)

	// Validate that we got smaller number of entities by 1
	_ = readEntries(ctx, t, store, query, len(tupdates)-1)

	assert.NoError(t, es.Purge(ctx))
	query = []*p4api.Entity{
		{Entity: &p4api.Entity_TableEntry{TableEntry: &p4api.TableEntry{}}},
	}
	readEntries(ctx, t, store, query, 0)
}

func generateRandomUpdates(info *p4info.P4Info) []*p4api.Update {
	updates := make([]*p4api.Update, 0, 512)
	tl := int32(len(info.Tables))
	for i := 0; i < cap(updates); i++ {
		tableInfo := info.Tables[rand.Int31n(tl)]
		for tableInfo.Size < 128 || tableInfo.IsConstTable {
			tableInfo = info.Tables[rand.Int31n(tl)]
		}
		entry := testutils.GenerateTableEntry(tableInfo, rand.Int31n(10), nil)
		update := &p4api.Update{Type: p4api.Update_INSERT, Entity: &p4api.Entity{Entity: &p4api.Entity_TableEntry{TableEntry: entry}}}
		updates = append(updates, update)
	}
	return updates
}

func readEntries(ctx context.Context, t *testing.T, store EntityStore, query []*p4api.Entity, count int) []*p4api.Entity {
	ch := make(chan *p4api.Entity, 1024)
	errs := store.Read(ctx, query, ch)
	for _, er := range errs {
		assert.NoError(t, er)
	}

	// Validate that we got all entries
	entities := make([]*p4api.Entity, 0, count)
	for e := range ch {
		entities = append(entities, e)
	}
	assert.Len(t, entities, count)
	return entities
}
