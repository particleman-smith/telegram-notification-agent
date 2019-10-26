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
	var returnStr = "ZFS Event 'Test' received."
	fmt.Println(returnStr)

	fmt.Println("Sending Telegram message.")
	telegramBot.SendMessage("Hey, I just got a test message!")
	fmt.Println("Message sent.")

	json.NewEncoder(writer).Encode(returnStr)
}
