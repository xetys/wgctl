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
	"os"
	"errors"
	"encoding/json"
	"io/ioutil"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new folder and configuration for wgctl",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		pwd, err := os.Getwd()
		FatalOnError(err)

		dir := fmt.Sprintf("%s/%s", pwd, name)

		exist := true
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			exist = false
		}

		if exist {
			FatalOnError(errors.New("there is already a directory with this name"))
		}

		err = os.Mkdir(dir, 0755)
		FatalOnError(err)

		config := &Config{Name: name}
		configBytes, err := json.Marshal(config)
		FatalOnError(err)

		err = ioutil.WriteFile(dir + "/config.json", configBytes, 0755)
		FatalOnError(err)

		fmt.Printf("wireguard configuration %s created", name)
		fmt.Println()
		fmt.Printf("run 'cd %s' to use the wgctl commands", name)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
