// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/onosproject/onos-net-lib/pkg/p4utils"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

// PipelineTranslator is an abstraction of an entity capable of translating high-level pipeline
// entities into low-level pipeline ones.
type PipelineTranslator interface {
	// Translate translates the given high-level pipeline entities into low-level pipeline ones.
	Translate(entities *[]p4api.Entity) *[]p4api.Entity

	// FromPipeline returns the P4 information describing the high-level pipeline
	FromPipeline() *p4info.P4Info

	// ToPipeline returns the P4 information describing the low-level target pipeline
	ToPipeline() *p4info.P4Info

	// For now, let's assume that the high-level and low-level pipelines use the same meta-data definitions.
	// If this doesn't hold, we will need to provide a mechanism for transcoding one into the other here.
}

// Provides identity pipeline entity translation
type identityTranslator struct {
	PipelineTranslator
	p4info *p4info.P4Info
}

// NewIdentityTranslator returns a new identity pipeline translator for the specified pipeline info
func NewIdentityTranslator(info *p4info.P4Info) PipelineTranslator {
	return &identityTranslator{p4info: info}
}

// NewIdentityTranslatorFromFile returns a new identity pipeline translator for pipeline info loaded from the given file
func NewIdentityTranslatorFromFile(path string) (PipelineTranslator, error) {
	info, err := p4utils.LoadP4Info(path)
	if err != nil {
		return nil, err
	}
	return NewIdentityTranslator(info), nil
}

// Translate returns the same entities as what was provided to it.
func (t *identityTranslator) Translate(entities *[]p4api.Entity) *[]p4api.Entity {
	return entities
}

// FromPipeline returns the P4 information describing the high-level pipeline; same as target pipeline
func (t *identityTranslator) FromPipeline() *p4info.P4Info {
	return t.p4info
}

// ToPipeline returns the P4 information describing the low-level target pipeline; same as high-level pipeline
func (t *identityTranslator) ToPipeline() *p4info.P4Info {
	return t.p4info
}
