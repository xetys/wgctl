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

// clientAddCmd represents the clientAdd command
var clientAddCmd = &cobra.Command{
	Use:   "add NAME ADDRESS ADDRESS_SPACE",
	Short: "Adds a new remote client",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		return nil
	},
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clientAdd called")
		name := args[0]
		address := args[1]
		cidr := args[2]
		newClient := Client{
			Name:    name,
			Address: address,
			CIDR:    cidr,
		}

		for idx, client := range wgContext.Config.Clients {
			if client.Name == name {
				wgContext.Config.Clients[idx] = newClient

				wgContext.SaveConfig()
				fmt.Println("client updated")
				return
			}
		}

		wgContext.Config.Clients = append(wgContext.Config.Clients, newClient)
		wgContext.SaveConfig()
		fmt.Println("client added")
	},
}

func init() {
	clientCmd.AddCommand(clientAddCmd)
}
