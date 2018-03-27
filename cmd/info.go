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
		fmt.Println("SSH KEYS: ")
		for _, sshKey := range wgContext.Config.SSHKeys {
			fmt.Println("- " + sshKey.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
