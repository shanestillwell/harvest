// Copyright ©2016 Ben Tranter <ben.tranter@metalabdesign.com>
//
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
//
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker-credential-helpers/osxkeychain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretStore = &osxkeychain.Osxkeychain{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Use Harvest from your terminal.",
	Long: `A command line interface for Harvest - the fast and simple way to schedule
your team.`,
	// Uncomment the following line if your bare application has an action associated with it
	//	Run: func(cmd *cobra.Command, args []string) { },
}

//Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// Read in config file and ENV variables if set.
func initConfig() {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".harvest")

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if _, err := os.Stat(cfgFile); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Reading initialization failed: %s\n", err.Error())
		}
	}
}
