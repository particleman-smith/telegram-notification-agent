/*
Controls sending and receipt of messages through Telegram using the
telegramsender and telegramreceiver
*/

package telegram

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	BotAPI   tgbotapi.BotAPI
	APIToken string
	Owner    int64
}

/*
Secrets are used for unmarshaling the data from secrets.json
*/
type Secrets struct {
	BotAPIToken string `json:"BotAPIToken`
	SudoerID    int64  `json:"RecipientID`
}

/*
NewBot creates a new instance of the Bot struct
*/
func NewBot() *Bot {
	println("Creating Telegram Api bot")

	b := new(Bot)

	// Try to get the API token from the env
	b.APIToken = os.Getenv("bot_token")
	if b.APIToken == "" {
		println("Using secrets.json for config")
		// Get APIToken from secrets.json
		jsonFile, jsonErr := os.Open("./secrets.json")
		if jsonErr != nil {
			log.Panic(jsonErr)
		}

		defer jsonFile.Close()

		// Read opened as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var secrets Secrets
		// Unmarshal JSON into secrets
		json.Unmarshal(byteValue, &secrets)
		b.APIToken = secrets.BotAPIToken
		b.Owner = secrets.SudoerID
	} else {
		println("Found environment variables")
		// Convert to int64
		owner, err := strconv.ParseInt(os.Getenv("sudoer_id"), 10, 64)
		if err == nil {
			b.Owner = owner
		}
	}

	setupErr := b.CheckRequirements()

	if setupErr != nil {
		log.Panic(setupErr)
	}

	bot, err := tgbotapi.NewBotAPI(b.APIToken)
	if err != nil {
		log.Panic(err)
	}

	b.BotAPI = *bot

	defer b.ListenForUpdates()

	return b
}

/*
SendMessage sends a message via Telegram
*/
func (bot Bot) SendMessage(body string) error {
	msg := tgbotapi.NewMessage(bot.Owner, body)
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

func (b Bot) CheckRequirements() error {
	if b.APIToken == "" || b.Owner <= 0 {
		return errors.New("Either the APIToken, Owner, or both were not found for this bot. Please check the secrets.json or env.conf files")
	}
	return nil
}
