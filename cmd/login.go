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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	// these need to be behind build flags for actual cross platform support
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
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
	loginCmd.PersistentFlags().String("org", "", "the organization you belong to on Harvest")
}

// login accepts a username and password, checking to make sure that they match
// a valid Harvest account. If they do, the username and password are stored in
// the system's secure vault (currently only keychain on OS X).
func login(cmd *cobra.Command, args []string) {
	id, err := cmd.Flags().GetString("email")
	if err != nil {
		fmt.Printf("Couldn't read the value for email: %s\n", err.Error())
		os.Exit(-1)
	}

	org, err := cmd.Flags().GetString("org")
	if err != nil {
		fmt.Printf("Couldn't read the value for org: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Please enter the password for your Harvest account:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Printf("Error while reading your password: %s\n", err.Error())
	}
	secret := string(password)

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
	defer resp.Body.Close()

	u := &struct {
		User *user `json:"user"`
	}{}
	json.NewDecoder(resp.Body).Decode(&u)

	switch resp.StatusCode {
	case 401:
		fmt.Printf("Incorrect email or password, pleas try again.\n")
		os.Exit(-1)
	case 404:
		fmt.Println(`Organization "` + org + `" could not be found, please check the value of the --org flag.`)
		os.Exit(-1)
	case 200:
		c := &credentials.Credentials{
			Username:  id,
			Secret:    cred, // use the base64 encoded string since it's what we pass to the client
			ServerURL: "https://" + org + ".harvestapp.com",
		}
		err := secretStore.Add(c)
		// thanks Docker for not exporting your error type
		if err.Error() == "The specified item already exists in the keychain." {
			fmt.Printf("You're already logged in!\n")
			return
		}
		if err != nil {
			fmt.Printf("Couldn't save credentials to Keychain, %s.\n", err.Error())
			os.Exit(-1)
		}
		fmt.Printf("Login successful, welcome to Harvest!\n")

		viper.Set("org", c.ServerURL)
		viper.Set("user_id", u.User.ID)
	default:
		fmt.Printf("Unknown error with status code %d. Could you file a bug to github.com/bentranter/harvest?\n", resp.StatusCode)
	}
}
