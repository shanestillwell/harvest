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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// these need to be behind build flags for actual cross platform support
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd respresents the login command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add your hours for today.",
	Long:  `Add your hours with the default being today, or specify a day. You can also add a message.`,
	Run:   add,
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().String("message", "", "the message to be sent")
	addCmd.Flags().String("day", "", "the day you worked the given hours (must be in the form 'YYYY-MM-DD')")
}

// add accepts a username and password, checking to make sure that they match
// a valid Harvest account. If they do, the username and password are stored in
// the system's secure vault (currently only keychain on OS X).
func add(cmd *cobra.Command, args []string) {
	msg, err := cmd.Flags().GetString("message")
	if err != nil {
		fmt.Printf("Couldn't read the value for message: %s\n", err.Error())
		os.Exit(-1)
	}

	day, err := cmd.Flags().GetString("day")
	if err != nil {
		fmt.Printf("Couldn't read the value for day: %s\n", err.Error())
		os.Exit(-1)
	}

	org := viper.GetString("org")
	cred, _, err := secretStore.Get(org)
	if err != nil {
		fmt.Printf("Looks like you haven't signed in yet. Run the command 'harvest login' to sign in. Additional info: %s\n", err.Error())
		os.Exit(-1)
	}

	data, err := json.M

	req, err := http.NewRequest("POST", "https://"+org+".harvestapp.com/account/who_am_i", nil)
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
	default:
		fmt.Printf("Unknown error with status code %d. Could you file a bug to github.com/bentranter/harvest?\n", resp.StatusCode)
	}
}
