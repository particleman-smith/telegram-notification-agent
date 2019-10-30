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

	defer b.ListenForUpdates()

	return b
}

/*
SendMessage sends a message via Telegram
*/
func (bot Bot) SendMessage(body string) error {
	msg := tgbotapi.NewMessage(bot.Secrets.RecipientID, body)
	message, err := bot.BotAPI.Send(msg)

	if err != nil {
		println("Error while sending message " + message.Text)
		log.Fatal(err)
		return err
	}

	return nil
}

/*
ListenForUpdates is a wrapper for the ListenForUpdatesRoutine
*/
func (bot Bot) ListenForUpdates() {
	go bot.ListenForUpdatesRoutine()
}

/*
ListenForUpdatesRoutine intemittently receives updates and responds to commands
*/
func (bot Bot) ListenForUpdatesRoutine() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.BotAPI.GetUpdatesChan(u)

	if err != nil {
		println("Error getting updates.")
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		ReplyToCommand(&bot.BotAPI, update.Message)
	}
}

func ReplyToCommand(botAPI *tgbotapi.BotAPI, incoming *tgbotapi.Message) {
	if incoming.IsCommand() {
		msg := tgbotapi.NewMessage(incoming.Chat.ID, "")
		switch incoming.Command() {
		case "start":
			msg.Text = "I am not good. Nor am I evil. I am no hero. Nor am I villain. I am AIDAN."
		case "help":
			msg.Text = "type /sayhi or /status."
		case "sayhi":
			msg.Text = "Hello."
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command."
		}
		botAPI.Send(msg)
	}
}
