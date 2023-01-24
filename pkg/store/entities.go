// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"crypto/sha1"
	_map "github.com/atomix/go-sdk/pkg/primitive/map"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
	"hash"
	"sort"
)

/*
Notes:

Create a map per device, per:
- each table (map[tableId]*table; map[hash(entry.matches)]*entry)
- all meters (map[meterId]*entry; same for all below)
- all counters
- all clone session entries
- all multicast group entries
- all action profile members
- all action profile groups

Also, later:
- all value set entries
- all register entries
- all extern entries
- all digest entries

Periodic reconciliation should result in updating the counters; direct or otherwise
Read requests for direct meters and counters should be done straight from the device (and update the stores)

*/

type table struct {
	info    *p4info.Table
	entries _map.Map[string, *p4api.TableEntry]
}

func (s *entityStore) modifyTableEntry(ctx context.Context, entry *p4api.TableEntry, insert bool) error {
	t, key, err := s.findTableAndKey(entry)
	if err != nil {
		return err
	}

	if insert {
		_, err = t.entries.Insert(ctx, key, entry)
	} else {
		_, err = t.entries.Update(ctx, key, entry)
	}
	if err != nil {
		return errors.FromAtomix(err)
	}

	// TODO: Implement updates of direct resources, if/as necessary
	return nil
}

func (s *entityStore) findTableAndKey(entry *p4api.TableEntry) (*table, string, error) {
	t, ok := s.tables[entry.TableId]
	if !ok {
		return nil, "", errors.NewInvalid("No such table %d", entry.TableId)
	}

	// Order field matches in canonical order based on field ID
	sortFieldMatches(entry.Match)

	// Produce a hash of the priority and the field matches to serve as a key
	key, err := t.entryKey(entry)
	if err != nil {
		return nil, "", err
	}
	return t, key, nil
}

// Produces a table entry key using a uint64 hash of its field matches; returns error if the matches do not comply
// with the table schema
func (t *table) entryKey(entry *p4api.TableEntry) (string, error) {
	if entry.IsDefaultAction {
		if len(entry.Match) > 0 {
			return "", errors.NewInvalid("Default action entry cannot have any match fields")
		}
		return "default", nil
	}

	hf := sha1.New()

	// This assumes matches have already been put in canonical order
	// TODO: implement field ID validation against the P4Info table schema
	for i, m := range entry.Match {
		if err := t.validateMatch(i, m); err != nil {
			return "", err
		}
		switch {
		case m.GetExact() != nil:
			_, _ = hf.Write([]byte{0x01})
			_, _ = hf.Write(m.GetExact().Value)
		case m.GetLpm() != nil:
			_, _ = hf.Write([]byte{0x02})
			writeHash(hf, m.GetLpm().PrefixLen)
			_, _ = hf.Write(m.GetLpm().Value)
		case m.GetRange() != nil:
			_, _ = hf.Write([]byte{0x03})
			_, _ = hf.Write(m.GetRange().Low)
			_, _ = hf.Write(m.GetRange().High)
		case m.GetTernary() != nil:
			_, _ = hf.Write([]byte{0x04})
			_, _ = hf.Write(m.GetTernary().Mask)
			_, _ = hf.Write(m.GetTernary().Value)
		case m.GetOptional() != nil:
			_, _ = hf.Write([]byte{0x05})
			_, _ = hf.Write(m.GetOptional().Value)
		}
	}
	return string(hf.Sum(nil)), nil
}

// Validates that the specified match corresponds to the expected table schema
func (t *table) validateMatch(i int, m *p4api.FieldMatch) error {
	if i >= len(t.info.MatchFields) {
		return errors.NewInvalid("Unexpected field match %d: %v", i, m)
	}

	// TODO: implement validation that the match is of expected type
	return nil
}

func writeHash(hash hash.Hash, n int32) {
	_, _ = hash.Write([]byte{byte((n & 0xff0000) >> 24), byte((n & 0xff0000) >> 16), byte((n & 0xff00) >> 8), byte(n & 0xff)})
}

// Sorts the given array of field matches in place based on the field ID
func sortFieldMatches(matches []*p4api.FieldMatch) {
	if len(matches) > 0 {
		sort.SliceStable(matches, func(i, j int) bool { return matches[i].FieldId < matches[j].FieldId })
	}
}

func (s *entityStore) removeTableEntry(ctx context.Context, entry *p4api.TableEntry) error {
	t, key, err := s.findTableAndKey(entry)
	if err != nil {
		return err
	}

	_, err = t.entries.Remove(ctx, key)
	if err != nil {
		return errors.FromAtomix(err)
	}
	return nil
}

func (s *entityStore) modifyCounterEntry(ctx context.Context, entry *p4api.CounterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyDirectCounterEntry(ctx context.Context, entry *p4api.DirectCounterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyMeterEntry(ctx context.Context, entry *p4api.MeterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyDirectMeterEntry(ctx context.Context, entry *p4api.DirectMeterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyActionProfileGroup(ctx context.Context, group *p4api.ActionProfileGroup, insert bool) error {
	return nil
}

func (s *entityStore) deleteActionProfileGroup(ctx context.Context, group *p4api.ActionProfileGroup) error {
	return nil
}

func (s *entityStore) modifyActionProfileMember(ctx context.Context, member *p4api.ActionProfileMember, insert bool) error {
	return nil
}

func (s *entityStore) deleteActionProfileMember(ctx context.Context, member *p4api.ActionProfileMember) error {
	return nil
}

func (s *entityStore) modifyMulticastGroupEntry(ctx context.Context, entry *p4api.MulticastGroupEntry, insert bool) error {
	return nil
}

func (s *entityStore) deleteMulticastGroupEntry(ctx context.Context, entry *p4api.MulticastGroupEntry) error {
	return nil
}

func (s *entityStore) modifyCloneSessionEntry(ctx context.Context, entry *p4api.CloneSessionEntry, insert bool) error {
	return nil
}

func (s *entityStore) deleteCloneSessionEntry(ctx context.Context, entry *p4api.CloneSessionEntry) error {
	return nil
}

func (s *entityStore) readTableEntries(ctx context.Context, entry *p4api.TableEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readCounterEntries(ctx context.Context, entry *p4api.CounterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readDirectCounterEntries(ctx context.Context, entry *p4api.DirectCounterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readMeterEntries(ctx context.Context, entry *p4api.MeterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readDirectMeterEntries(ctx context.Context, entry *p4api.DirectMeterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readActionProfileGroups(ctx context.Context, group *p4api.ActionProfileGroup, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readActionProfileMembers(ctx context.Context, member *p4api.ActionProfileMember, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readMulticastGroupEntries(ctx context.Context, entry *p4api.MulticastGroupEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readCloneSessionEntries(ctx context.Context, entry *p4api.CloneSessionEntry, ch chan<- p4api.Entity) error {
	return nil
}
