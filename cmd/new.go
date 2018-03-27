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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
