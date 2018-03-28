// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"errors"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Executes the config to your nodes. Add -d or --dry-run if you want to check the configs before",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		if len(wgContext.Config.Nodes) == 0 {
			return errors.New("you must add at least one node")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		ifaceNumber, _ := cmd.Flags().GetInt("iface-number")
		wgContext.Config.SetupEncryptedNetwork(ifaceNumber, dryRun)
		wgContext.SaveConfig()
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.Flags().BoolP("dry-run", "d", false, "if true, wgctl only prints the distinct configs")
	doCmd.Flags().IntP("iface-number", "i", 0, "number of the interface. Will generate a wg0.conf, wg1.conf, ...")
}
