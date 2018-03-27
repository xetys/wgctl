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
	"text/tabwriter"
	"os"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints out the current configuration",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !wgContext.IsConfigLoaded() {
			return errors.New("you must be inside a config directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NAME: %s\n", wgContext.Config.Name)
		fmt.Println()
		fmt.Println("NODES: ")
		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(tw, "NAME\tEXTERNAL IP\tINTERNAL IP\tSSH KEY\t")

		for _, node := range wgContext.Config.Nodes {
			fmt.Fprintf(tw,"%s\t%s\t%s\t%s\t", node.Name, node.IPAddress, node.PrivateIPAddress, node.SSHKeyName)
			fmt.Fprintln(tw)
		}
		tw.Flush()
		fmt.Println()

		fmt.Println("CLIENTS: ")
		tw.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(tw, "NAME\tADDRESS\tADDRESS SPACE\t")
		for _, client := range wgContext.Config.Clients {
			fmt.Fprintf(tw,"%s\t%s\t%s\t", client.Name, client.Address, client.CIDR)
			fmt.Fprintln(tw)
		}
		tw.Flush()
		fmt.Println()

		fmt.Println("SSH KEYS: ")
		for _, sshKey := range wgContext.Config.SSHKeys {
			fmt.Println("- " + sshKey.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
