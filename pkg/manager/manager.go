// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package manager is is the main coordinator for the ONOS control subsystem.
package manager

import (
	log "k8s.io/klog"
)

var mgr Manager

// Manager single point of entry for the config system.
type Manager struct {
}

// NewManager initializes the network control manager subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")
	mgr = Manager{}

	return &mgr, nil
}

// LoadManager creates a configuration subsystem manager primed with stores loaded from the specified files.
func LoadManager() (*Manager, error) {
	return NewManager()
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
	// Start the main dispatcher system
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
