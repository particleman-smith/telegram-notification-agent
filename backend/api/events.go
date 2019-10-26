/*
Responsd to ZFS event HTTP requests
*/

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/particleman-smith/telegram-notification-agent/backend/telegram"
)

var telegramBot = telegram.NewBot()

/*
Test sends the given message via the Telegram Bot
*/
func Test(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("ZFS Event 'Test' received.")

	fmt.Println("Sending Telegram message.")
	err := telegramBot.SendMessage("Hey, I just got a test message!")

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
	err := error(nil)
	switch request.URL.Path {
	case "backup-event/failure":
		err = telegramBot.SendMessage("ERROR\nI encountered a failure backing up /home, /etc, or /var.")
		break
	}

	if err != nil {
		returnMsg := "Message failed to send."
		fmt.Println(returnMsg)
		json.NewEncoder(writer).Encode(returnMsg)
	}
}
