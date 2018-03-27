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
	"github.com/xetys/hetzner-kube/pkg/clustermanager"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		if len(args) != 4 {
			return errors.New(fmt.Sprintf("you must pass 4 arguments. %d passed", len(args)))
		}

		sshKeyName := args[3]

		for _, key := range wgContext.Config.SSHKeys {
			if key.Name == sshKeyName {
				return nil
			}
		}

		return errors.New("key not found")
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		exIp := args[1]
		inIp := args[2]
		sshKeyName := args[3]

		newNode := clustermanager.Node{
			Name: name,
			IPAddress: exIp,
			PrivateIPAddress: inIp,
			SSHKeyName: sshKeyName,
		}

		found := false
		for idx, node := range wgContext.Config.Nodes {
			if node.Name == newNode.Name {
				wgContext.Config.Nodes[idx] = newNode
				found = true
				break
			}
		}

		if !found {
			wgContext.Config.Nodes = append(wgContext.Config.Nodes, newNode)
		}

		wgContext.SaveConfig()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
