/*
Controller methods to respond to event HTTP requests
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var _secrets = GetSecrets()
var _telegramBot = NewBot(*_secrets)

/*
Test sends the given message via the Telegram Bot
*/
func Test(writer http.ResponseWriter, request *http.Request) {
	if !CheckAPIToken(request.Header.Get("Access-Token")) {
		http.Error(writer, "Invalid access token", 401)
		return
	}

	fmt.Println("ZFS Event 'Test' received.")

	fmt.Println("Sending Telegram message.")
	err := _telegramBot.SendMessage("I just received a test API message.")

	returnMsg := ""

	if err != nil {
		returnMsg = "Message failed to send."

		fmt.Println(err.Error())
	} else {
		returnMsg = "Message sent."
	}

	fmt.Println(returnMsg)
	json.NewEncoder(writer).Encode(returnMsg)
}

/*
Error sends an error message via the Telegram Bot. The error message is based on the request URL path.
*/
func Error(writer http.ResponseWriter, request *http.Request) {
	if !CheckAPIToken(request.Header.Get("Access-Token")) {
		http.Error(writer, "Invalid access token", 401)
		return
	}

	bodyMap, readErr := GetBodyMapFromRequest(request)

	if readErr != nil {
		http.Error(writer, "Could not parse body.\n"+readErr.Error(), 500)
		return
	}

	send := true // Whether or not to send the message (body contents may warrant not sending a message)
	msg := "WARNING\n"

	// Error types
	switch request.URL.Path {
	// Bash
	case "/bash-event/exec-failure":
		script := bodyMap["script"]
		reason := bodyMap["reason"]
		msg += "I encountered a failure running the Python script `" + fmt.Sprintf("%v", script) + "`."
		msg += "\n" + fmt.Sprintf("%v", reason)
	// ZFS
	case "/zfs-event/data-error":
		msg += "I noticed data corruption on a ZFS vdev."
	case "/zfs-event/zpool-state":
		state := fmt.Sprintf("%v", bodyMap["status"])
		msg += "A a ZFS vdev has entered " + AOrAn(state) + state + " state!"
	// Backup
	case "/backup-event/failure":
		msg += "I encountered a failure backing up /home, /etc, or /var."
	}

	if send {
		// Send the message via Telegram
		sendErr := _telegramBot.SendMessage(msg)

		// Handle Telegram send errors
		if sendErr != nil {
			returnMsg := "Message failed to send."
			fmt.Println(returnMsg)
			json.NewEncoder(writer).Encode(returnMsg)
		}
	}
}

/*
Info sends a message containing misc info via the the Telegram Bot. The message is based on the request URL path.
*/
func Info(writer http.ResponseWriter, request *http.Request) {
	if !CheckAPIToken(request.Header.Get("Access-Token")) {
		http.Error(writer, "Invalid access token", 401)
		return
	}

	_, err := GetBodyMapFromRequest(request)

	if err != nil {
		http.Error(writer, "Could not parse body.\n"+err.Error(), 500)
		return
	}

	send := true
	msg := "INFO\n"

	// Info types
	switch request.URL.Path {
	case "/backup-event/success":
		msg += "AIDAN backed up successfully."
	}

	if send {
		// Send the message via Telegram
		sendErr := _telegramBot.SendMessage(msg)

		// Handle Telegram send errors
		if sendErr != nil {
			returnMsg := "Message failed to send."
			fmt.Println(returnMsg)
			json.NewEncoder(writer).Encode(returnMsg)
		}
	}
}

/*
CheckAPIToken compares the given secret the the APIToken in the config and returns whether or not they match
*/
func CheckAPIToken(secret string) bool {
	if secret == _secrets.AccessToken {
		return true
	}
	return false
}

/*
GetBodyMapFromRequest reads the body of an http request and parses it into a map[string]interface{}
*/
func GetBodyMapFromRequest(request *http.Request) (map[string]interface{}, error) {
	// Read body and handle errors
	body, readErr := ioutil.ReadAll(request.Body)
	if readErr != nil {
		return nil, readErr
	}
	defer request.Body.Close()

	// Parse the body into an interface
	var bodyMap map[string]interface{}
	parseErr := json.Unmarshal([]byte(body), &bodyMap)

	if parseErr != nil {
		return nil, parseErr
	}

	return bodyMap, nil
}

/*
IsFirstCharVowel checks if the first character of a word is a vowel and returns true or false.
*/
func IsFirstCharVowel(word string) bool {
	vowels := []byte{'a', 'A', 'e', 'E', 'i', 'I', 'o', 'O', 'u', 'U'}
	for _, v := range vowels {
		if word[0] == v {
			return true
		}
	}
	return false
}

/*
AOrAn will "a" or "an" based on whether or not the first character of the given word is a vowel.
*/
func AOrAn(word string) string {
	if IsFirstCharVowel(word) {
		return "an "
	}
	return "a "
}
