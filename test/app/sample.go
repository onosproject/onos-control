// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// Package app hosts code snippets simulating various application use-cases of the onos-control library
package app

import (
	"context"
	"github.com/atomix/go-sdk/pkg/client"
	"github.com/onosproject/onos-control/pkg/api"
	"github.com/onosproject/onos-control/pkg/controller"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-net-lib/pkg/p4utils"
)

var log = logging.GetLogger("sample")

// InitExample sketches out steps necessary to initialize the library, add a device controller and query its status.
func InitExample() error {
	// Node: Consider creating a StratumRoleBuilder
	role := p4utils.NewStratumRole("sample", 0, []byte{}, false, false)
	devices := controller.NewController(role, client.NewClient())

	// Note: Consider including a utility to easily add all devices from onos-topo using a realm label.
	// This would fetch all devices matching the realm-label and the required onos.topo.StratumAgents, and
	// onos.provisioner.DeviceConfig aspects, which it would use to get the P4Info from the device provisioner and for
	// P4Runtime endpoint. This would require either specifying or internally creating onos-topo and device-provisioner
	// client. E.g.  sdk.AddRealmDevices(ctx, devices, topoClient, provisionerClient, realmLabel, realmValue)

	// Create a translator from the specified p4info file
	translator, err := api.NewIdentityTranslatorFromFile("../../test/p4info.txt")
	if err != nil {
		log.Errorf("Unable to create pipeline translator: %+v", err)
		return err
	}

	ctx := context.Background()
	fooDevice, err := devices.Add(ctx, "foo", "fabric-sim:20000", translator)
	if err != nil {
		log.Errorf("Unable to add controller for device %s: %+v", "foo", err)
		return err
	}

	log.Infof("Controller state: %s", fooDevice.State())
	return nil
}
