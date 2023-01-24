// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

// EntityStore is an abstraction of a store capable of maintaining various P4 entities
type EntityStore interface {
	// ID returns the ID of the device whose control entities it persists.
	ID() topo.ID

	// P4Info returns the P4Info used to structure the store and validate entries
	P4Info() *p4info.P4Info

	// Read accepts a query in form of a list of partially populated entities and returns any
	// matching entities on the specified channel
	Read(ctx context.Context, query []*p4api.Entity, ch chan<- p4api.Entity) []error

	// Write persists the specified list of updates.
	Write(ctx context.Context, updates []*p4api.Update) error
}

type entityStore struct {
	EntityStore
	id   topo.ID
	info *p4info.P4Info
	// TODO: Insert Atomix primitives to track table, group, meter, etc. entries
}

// Creates a new P4 entity store for the specified device
func newEntityStore(id topo.ID, info *p4info.P4Info) (EntityStore, error) {
	// TODO: Implement me
	return &entityStore{id: id, info: info}, nil
}

// ID returns the ID of the device whose control entities it persists.
func (s *entityStore) ID() topo.ID {
	return s.id
}

// Read accepts a query in form of a list of partially populated entities and returns any
// matching entities on the specified channel
func (s *entityStore) Read(ctx context.Context, query []*p4api.Entity, ch chan<- p4api.Entity) []error {
	// Allocate the same number of errors as there are requests - expressed as entities
	errors := make([]error, len(query))

	for i, request := range query {
		errors[i] = s.processRead(request, ch)
	}
	return errors
}

func (s *entityStore) processRead(query *p4api.Entity, ch chan<- p4api.Entity) error {
	switch {
	case query.GetTableEntry() != nil:
		return s.readTableEntries(query.GetTableEntry(), ch)
	case query.GetCounterEntry() != nil:
		return s.readCounterEntries(query.GetCounterEntry(), ch)
	case query.GetDirectCounterEntry() != nil:
		return s.readDirectCounterEntries(query.GetDirectCounterEntry(), ch)
	case query.GetMeterEntry() != nil:
		return s.readMeterEntries(query.GetMeterEntry(), ch)
	case query.GetDirectMeterEntry() != nil:
		return s.readDirectMeterEntries(query.GetDirectMeterEntry(), ch)

	case query.GetActionProfileGroup() != nil:
		return s.readActionProfileGroups(query.GetActionProfileGroup(), ch)
	case query.GetActionProfileMember() != nil:
		return s.readActionProfileMembers(query.GetActionProfileMember(), ch)

	case query.GetPacketReplicationEngineEntry() != nil:
		switch {
		case query.GetPacketReplicationEngineEntry().GetMulticastGroupEntry() != nil:
			return s.readMulticastGroupEntries(query.GetPacketReplicationEngineEntry().GetMulticastGroupEntry(), ch)
		case query.GetPacketReplicationEngineEntry().GetCloneSessionEntry() != nil:
			return s.readCloneSessionEntries(query.GetPacketReplicationEngineEntry().GetCloneSessionEntry(), ch)
		}

	case query.GetRegisterEntry() != nil:
	case query.GetValueSetEntry() != nil:
	case query.GetDigestEntry() != nil:
	case query.GetExternEntry() != nil:
	default:
	}
	return nil
}

// Write persists the specified list of updates.
func (s *entityStore) Write(ctx context.Context, updates []*p4api.Update) error {
	for _, update := range updates {
		switch {
		case update.Type == p4api.Update_INSERT:
			if err := s.processModify(update, true); err != nil {
				log.Warnf("Device %s: Unable to insert entry: %+v", s.id, err)
				return err
			}
		case update.Type == p4api.Update_MODIFY:
			if err := s.processModify(update, false); err != nil {
				log.Warnf("Device %s: Unable to update entry: %+v", s.id, err)
				return err
			}
		case update.Type == p4api.Update_DELETE:
			if err := s.processDelete(update); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *entityStore) processModify(update *p4api.Update, isInsert bool) error {
	entity := update.Entity
	var err error
	switch {
	case entity.GetTableEntry() != nil:
		err = s.modifyTableEntry(entity.GetTableEntry(), isInsert)
	case entity.GetCounterEntry() != nil:
		err = s.modifyCounterEntry(entity.GetCounterEntry(), isInsert)
	case entity.GetDirectCounterEntry() != nil:
		err = s.modifyDirectCounterEntry(entity.GetDirectCounterEntry(), isInsert)
	case entity.GetMeterEntry() != nil:
		err = s.modifyMeterEntry(entity.GetMeterEntry(), isInsert)
	case entity.GetDirectMeterEntry() != nil:
		err = s.modifyDirectMeterEntry(entity.GetDirectMeterEntry(), isInsert)

	case entity.GetActionProfileGroup() != nil:
		err = s.modifyActionProfileGroup(entity.GetActionProfileGroup(), isInsert)
	case entity.GetActionProfileMember() != nil:
		err = s.modifyActionProfileMember(entity.GetActionProfileMember(), isInsert)

	case entity.GetPacketReplicationEngineEntry() != nil:
		switch {
		case entity.GetPacketReplicationEngineEntry().GetMulticastGroupEntry() != nil:
			err = s.modifyMulticastGroupEntry(entity.GetPacketReplicationEngineEntry().GetMulticastGroupEntry(), isInsert)
		case entity.GetPacketReplicationEngineEntry().GetCloneSessionEntry() != nil:
			err = s.modifyCloneSessionEntry(entity.GetPacketReplicationEngineEntry().GetCloneSessionEntry(), isInsert)
		}

	case entity.GetRegisterEntry() != nil:
		log.Warnf("Device %s: RegisterEntry write is not supported yet: %+v", s.id, entity.GetRegisterEntry())
	case entity.GetValueSetEntry() != nil:
		log.Warnf("Device %s: ValueSetEntry write is not supported yet: %+v", s.id, entity.GetValueSetEntry())
	case entity.GetDigestEntry() != nil:
		log.Warnf("Device %s: DigestEntry write is not supported yet: %+v", s.id, entity.GetDigestEntry())
	case entity.GetExternEntry() != nil:
		log.Warnf("Device %s: ExternEntry write is not supported yet: %+v", s.id, entity.GetExternEntry())
	default:
	}
	return err
}

func (s *entityStore) processDelete(update *p4api.Update) error {
	entity := update.Entity
	var err error
	switch {
	case entity.GetTableEntry() != nil:
		err = s.removeTableEntry(entity.GetTableEntry())
	case entity.GetCounterEntry() != nil:
		return errors.NewInvalid("counter cannot be deleted")
	case entity.GetDirectCounterEntry() != nil:
		err = errors.NewInvalid("direct counter entry cannot be deleted")
	case entity.GetMeterEntry() != nil:
		return errors.NewInvalid("meter cannot be deleted")
	case entity.GetDirectMeterEntry() != nil:
		err = errors.NewInvalid("direct meter entry cannot be deleted")

	case entity.GetActionProfileGroup() != nil:
		err = s.deleteActionProfileGroup(entity.GetActionProfileGroup())
	case entity.GetActionProfileMember() != nil:
		err = s.deleteActionProfileMember(entity.GetActionProfileMember())

	case entity.GetPacketReplicationEngineEntry() != nil:
		switch {
		case entity.GetPacketReplicationEngineEntry().GetMulticastGroupEntry() != nil:
			err = s.deleteMulticastGroupEntry(entity.GetPacketReplicationEngineEntry().GetMulticastGroupEntry())
		case entity.GetPacketReplicationEngineEntry().GetCloneSessionEntry() != nil:
			err = s.deleteCloneSessionEntry(entity.GetPacketReplicationEngineEntry().GetCloneSessionEntry())
		}

	case entity.GetRegisterEntry() != nil:
	case entity.GetValueSetEntry() != nil:
	case entity.GetDigestEntry() != nil:
	case entity.GetExternEntry() != nil:
	default:
	}
	return err
}
