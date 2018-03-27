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
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var wgContext *WgCtlContext
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wgctl",
	Short: "A small tool for managing wireguard VPN",
	Long: `This tool enables management of remote wireguard setups.

Here is a quick explanation:

	# wgctl new my-servers
	-> this will create a new folder in the current directory and place a config.json with a name for usage. You must enter the folder for all further actions
	# cd my-servers
	# wgctl import-ssh-keys
	-> this will import your ssh keys for connecting the nodes
	# wgctl add 42.42.42.1 10.0.0.1 id_rsa
	# wgctl add 42.42.42.2 10.0.0.2 id_rsa
	# wgctl add 42.42.42.3 10.0.0.3 id_rsa
	-> this will add three nodes you are going to manage
	# wgctl install
	-> this step ensures your nodes will have wireguard install
	# wgctl do -d
	-> take a look into the configurations if you want
	# wgctl do
	-> setup the network

With these steps your servers should be connected. Now you could add an client config:
	# wgctl client add me 10.0.1.1 10.0.1.0/24
	-> adding a client 'me' to your list
	# wgctl do
	-> You must run this, to generate a certificate!
	# wgctl client gen-config me > wgme.conf
	-> generates a config for a client. You can share this one!
	# wg-quick up ./wgme.conf
	-> connect to the VPN
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wgctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	wgContext = &WgCtlContext{}

	_, err := wgContext.LoadFromLocalDir()

	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".wgctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".wgctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
