// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
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

func (s *entityStore) modifyTableEntry(entry *p4api.TableEntry, insert bool) error {
	return nil
}

func (s *entityStore) removeTableEntry(entry *p4api.TableEntry) error {
	return nil
}

func (s *entityStore) modifyCounterEntry(entry *p4api.CounterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyDirectCounterEntry(entry *p4api.DirectCounterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyMeterEntry(entry *p4api.MeterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyDirectMeterEntry(entry *p4api.DirectMeterEntry, insert bool) error {
	return nil
}

func (s *entityStore) modifyActionProfileGroup(group *p4api.ActionProfileGroup, insert bool) error {
	return nil
}

func (s *entityStore) deleteActionProfileGroup(group *p4api.ActionProfileGroup) error {
	return nil
}

func (s *entityStore) modifyActionProfileMember(member *p4api.ActionProfileMember, insert bool) error {
	return nil
}

func (s *entityStore) deleteActionProfileMember(member *p4api.ActionProfileMember) error {
	return nil
}

func (s *entityStore) modifyMulticastGroupEntry(entry *p4api.MulticastGroupEntry, insert bool) error {
	return nil
}

func (s *entityStore) deleteMulticastGroupEntry(entry *p4api.MulticastGroupEntry) error {
	return nil
}

func (s *entityStore) modifyCloneSessionEntry(entry *p4api.CloneSessionEntry, insert bool) error {
	return nil
}

func (s *entityStore) deleteCloneSessionEntry(entry *p4api.CloneSessionEntry) error {
	return nil
}

func (s *entityStore) readTableEntries(entry *p4api.TableEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readCounterEntries(entry *p4api.CounterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readDirectCounterEntries(entry *p4api.DirectCounterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readMeterEntries(entry *p4api.MeterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readDirectMeterEntries(entry *p4api.DirectMeterEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readActionProfileGroups(group *p4api.ActionProfileGroup, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readActionProfileMembers(member *p4api.ActionProfileMember, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readMulticastGroupEntries(entry *p4api.MulticastGroupEntry, ch chan<- p4api.Entity) error {
	return nil
}

func (s *entityStore) readCloneSessionEntries(entry *p4api.CloneSessionEntry, ch chan<- p4api.Entity) error {
	return nil
}
