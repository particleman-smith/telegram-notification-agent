/*
Controls retrieval and modeling of secrets from the secrets.json file (used while debugging)
or from environment variables
*/

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

/*
Secrets are used for unmarshaling the data from secrets.json
*/
type Secrets struct {
	BotAPIToken string `json:"BotAPIToken"` // The API token used to control the Telegram bot
	SudoerID    int64  `json:"SudoerID"`    // The owner of the bot, used as the recipient of messages triggered by events
	AccessToken string `json:"AccessToken"` // The secret that must be on the "Access_Token" header of any incoming HTTP request
}

/*
GetSecrets attempts to get the secrets from a secrets.json file or the environment variables
*/
func GetSecrets() *Secrets {
	s := new(Secrets)

	// Try to get the API token from the env
	s.BotAPIToken = os.Getenv("bot_token")
	if s.BotAPIToken == "" { // Get values from secrets.json
		println("Using secrets.json for config")

		jsonFile, jsonErr := os.Open("./secrets.json")
		CheckErr(jsonErr)

		defer jsonFile.Close()

		// Read opened as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Unmarshal JSON into secrets
		json.Unmarshal(byteValue, &s)
	} else {
		println("Found environment variables for config")
		err := error(nil)

		// BotAPIToken was already added above
		// Get SudoerID
		s.SudoerID, err = strconv.ParseInt(os.Getenv("sudoer_id"), 10, 64) // Convert to int64
		CheckErr(err)

		// Get AccessToken
		s.AccessToken = os.Getenv("access_token")
		if s.AccessToken == "" {
			err = errors.New("Could not get access_token from environment variables")
		}
		CheckErr(err)
	}

	return s
}

/*
CheckErr checks the given error and throws and panics if it is not nil
*/
func CheckErr(err error) {
	if err != nil {
		log.Println("Failed getting secrets:")
		log.Panic(err)
	}
}
