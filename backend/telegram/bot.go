/*
Controls sending and receipt of messages through Telegram using the
telegramsender and telegramreceiver
*/

package telegram

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	BotAPI  tgbotapi.BotAPI
	Secrets Secrets
}

/*
Secrets are used for unmarshaling the data from secrets.json
*/
type Secrets struct {
	BotAPIToken string `json:"BotAPIToken`
	RecipientID int64  `json:"RecipientID`
}

/*
NewBot creates a new instance of the Bot struct
*/
func NewBot() *Bot {
	// Get APIToken from secrets.json
	jsonFile, jsonErr := os.Open("./secrets.json")
	if jsonErr != nil {
		log.Panic(jsonErr)
	}

	// Read opened as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var secrets Secrets
	// Unmarshal JSON into secrets
	json.Unmarshal(byteValue, &secrets)

	b := new(Bot)
	b.Secrets = secrets
	bot, err := tgbotapi.NewBotAPI(b.Secrets.BotAPIToken)
	if err != nil {
		log.Panic(err)
	}

	defer jsonFile.Close()

	b.BotAPI = *bot
	return b
}

/*
SendMessage sends a message via Telegram
*/
func (bot Bot) SendMessage(body string) {
	msg := tgbotapi.NewMessage(bot.Secrets.RecipientID, body)
	bot.BotAPI.Send(msg)
}
