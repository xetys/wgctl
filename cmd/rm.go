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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a node from your list",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		return nil
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		for idx, node := range wgContext.Config.Nodes {
			if node.Name == name {
				wgContext.Config.Nodes = append(wgContext.Config.Nodes[:idx], wgContext.Config.Nodes[idx+1:]...)
				wgContext.SaveConfig()
				return
			}
		}

		fmt.Println("node not found")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
