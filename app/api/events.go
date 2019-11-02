/*
Responsd to ZFS event HTTP requests
*/

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/particleman-smith/telegram-notification-agent/app/telegram"
)

var telegramBot = telegram.NewBot()

/*
Test sends the given message via the Telegram Bot
*/
func Test(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("ZFS Event 'Test' received.")

	fmt.Println("Sending Telegram message.")
	err := telegramBot.SendMessage("I just received a test API message.")

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
	// Read body and handle errors
	body, readErr := ioutil.ReadAll(request.Body)
	if readErr != nil {
		http.Error(writer, readErr.Error(), 500)
		return
	}
	defer request.Body.Close()

	// Parse the body into an interface
	var bodyMap map[string]interface{}
	parseErr := json.Unmarshal([]byte(body), &bodyMap)

	if parseErr != nil {
		http.Error(writer, "Could not parse body.\n"+parseErr.Error(), 500)
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
		msg += "A disk in a ZFS vdev has entered " + AOrAn(state) + state + " state!"
	// Backup
	case "/backup-event/failure":
		msg += "I encountered a failure backing up /home, /etc, or /var."
	}

	if send {
		// Send the message via Telegram
		sendErr := telegramBot.SendMessage(msg)

		// Handle Telegram send errors
		if sendErr != nil {
			returnMsg := "Message failed to send."
			fmt.Println(returnMsg)
			json.NewEncoder(writer).Encode(returnMsg)
		}
	}
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
