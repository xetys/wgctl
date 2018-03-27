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
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"strings"
	"github.com/xetys/hetzner-kube/pkg/clustermanager"
)

// importSshKeysCmd represents the importSshKeys command
var importSshKeysCmd = &cobra.Command{
	Use:   "import-ssh-keys",
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
		fmt.Println("importSshKeys called")
		home, err := homedir.Dir()
		FatalOnError(err)
		sshPath := home + "/.ssh/"

		files, err := ioutil.ReadDir(sshPath)
		FatalOnError(err)
		wgContext.Config.SSHKeys = []clustermanager.SSHKey{}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileContent, err := ioutil.ReadFile(sshPath + file.Name())
			FatalOnError(err)
			if strings.Contains(string(fileContent), "-----BEGIN RSA PRIVATE KEY-----") {
				fmt.Printf("adding key '%s'", file.Name())
				fmt.Println()
				wgContext.Config.SSHKeys = append(wgContext.Config.SSHKeys, clustermanager.SSHKey{
					Name:           file.Name(),
					PrivateKeyPath: sshPath + file.Name(),
					PublicKeyPath:  sshPath + file.Name() + ".pub",
				})
			}
		}

		wgContext.SaveConfig()
	},
}

func init() {
	rootCmd.AddCommand(importSshKeysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importSshKeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importSshKeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
