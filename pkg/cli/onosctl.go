// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package cli holds ONOS command-line interface.
package cli

import (
	"fmt"
	"github.com/onosproject/onos-control/pkg/cli/command"
	"os"
)

// Execute runs the root command and any sub-commands.
func Execute() {
	rootCmd := command.GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
