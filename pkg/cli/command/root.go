// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package command holds ONOS command-line command implementations.
package command

import (
	"github.com/onosproject/onos-config/pkg/certs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetRootCommand returns the root CLI command.
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onos",
		Short: "ONOS command line client",
	}

	viper.SetDefault("address", ":5150")
	viper.SetDefault("keyPath", certs.Client1Key)
	viper.SetDefault("certPath", certs.Client1Crt)

	cmd.PersistentFlags().StringP("address", "a", viper.GetString("address"), "the controller address")
	cmd.PersistentFlags().StringP("keyPath", "k", viper.GetString("keyPath"), "path to client private key")
	cmd.PersistentFlags().StringP("certPath", "c", viper.GetString("certPath"), "path to client certificate")
	cmd.PersistentFlags().String("config", "", "config file (default: $HOME/.onos/config.yaml)")

	_ = viper.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag("keyPath", cmd.PersistentFlags().Lookup("keyPath"))
	_ = viper.BindPFlag("certPath", cmd.PersistentFlags().Lookup("certPath"))

	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newCompletionCommand())
	return cmd
}
