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
	"strings"
	"errors"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Given the target node is ubuntu based, it installs wireguard",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		for _, node := range wgContext.Config.Nodes {
			out, err := wgContext.SSH().RunCmd(node, "type -p wg > /dev/null &> /dev/null; echo $?")
			FatalOnError(err)
			out = strings.TrimSpace(out)
			installed := out == "0"

			if !installed {
				fmt.Printf("installing wireguard on node %s\n", node.Name)
				_, err = wgContext.SSH().RunCmd(node, "add-apt-repository ppa:wireguard/wireguard -y")
				FatalOnError(err)
				_, err = wgContext.SSH().RunCmd(node, "apt-get update && apt-get install -y wireguard linux-headers-$(uname -r) linux-headers-virtual")
				FatalOnError(err)
				fmt.Println("wireguard installed!")
			} else {
				fmt.Println("wireguard already exist")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
