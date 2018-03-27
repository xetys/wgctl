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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
