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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	// "github.com/docker/docker-credential-helpers/credentials"
	// "github.com/docker/docker-credential-helpers/osxkeychain"
	"github.com/spf13/cobra"
)

// loginCmd respresents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Harvest account via the CLI",
	Long:  `Login to your Harvest account by supplying a username and password.`,
	Run:   login,
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().String("email", "", "the email associated with your Harvest account")
	loginCmd.PersistentFlags().String("password", "", "the password associated with your Harvest account")
	loginCmd.PersistentFlags().String("org", "", "the organization you belong to on Harvest")
	// Here you will define your flags and configuration settings.
	// We need the --user, --email, --password, and --org flags

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )

}

// login accepts a username and password, checking to make sure that they match
// a valid Harvest account. If they do, the username and password are stored in
// the system's secure vault (currently only keychain on OS X).
func login(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("email")
	secret, _ := cmd.Flags().GetString("password")
	org, _ := cmd.Flags().GetString("org")
	cred := base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))

	req, err := http.NewRequest("GET", "https://"+org+".harvestapp.com/account/who_am_i", nil)
	if err != nil {
		fmt.Printf("Error creating new HTTP request: %s\n", err.Error())
		os.Exit(-1)
		return
	}
	req.Header.Set("Authorization", "Basic "+cred)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error during HTTP request: %s\n", err.Error())
		os.Exit(-1)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading who am i response: %s\n", err.Error())
	}
	defer resp.Body.Close()

	fmt.Printf("Response:\n\n%s\n", body)
}
