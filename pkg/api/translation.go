// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
)

// PipelineTranslator is an abstraction of an entity capable of translating high-level pipeline
// entities into low-level pipeline ones.
type PipelineTranslator interface {
	// Translate translates the given high-level pipeline entities into low-level pipeline ones.
	Translate(entities *[]p4api.Entity) *[]p4api.Entity
}

// IdentityTranslator provides identity pipeline entity translation
type IdentityTranslator struct {
	PipelineTranslator
}

// Translate returns the same entities as what was provided to it.
func (t *IdentityTranslator) Translate(entities *[]p4api.Entity) *[]p4api.Entity {
	return entities
}
