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
	"fmt"

	"github.com/spf13/cobra"
	"errors"
)

// genConfigCmd represents the genConfig command
var genConfigCmd = &cobra.Command{
	Use:   "gen-config",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		clientName := args[0]

		for _, client := range wgContext.Config.Clients {
			if client.Name == clientName {
				if client.KeyPair.Public == "" {
					fmt.Println("your client has no keys installed. Please run wgctl do to generate them")
					return
				}

				config := GenerateClientConf(client, wgContext.Config.Nodes)
				fmt.Println(config)
				return
			}
		}
	},
}

func init() {
	clientCmd.AddCommand(genConfigCmd)
}
