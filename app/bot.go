/*
Controls sending and receipt of messages through Telegram using the
telegramsender and telegramreceiver
*/

package main

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
Bot is used to store the Telegram BotAPI as well as its APIToke and the user ID of the owner.
*/
type Bot struct {
	BotAPI   tgbotapi.BotAPI // The interface to Telegram
	APIToken string          // The APIToken used to control the Telegram bot
	Owner    int64           // The user ID of the owner of the bot, used as the recipient of messages triggered by events
}

/*
NewBot creates a new instance of the Bot struct
*/
func NewBot(s Secrets) *Bot {
	println("Creating Telegram Api bot")

	b := new(Bot)

	b.APIToken = s.BotAPIToken
	b.Owner = s.SudoerID

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

/*
ReplyToCommand handles commands sent from a user to the bot
*/
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
			msg.Text = "Up and running."
		default:
			msg.Text = "I don't know that command."
		}
		botAPI.Send(msg)
	}
}

/*
CheckRequirements ensures that the necessary secrets for the given Bot have been obtained
*/
func (bot Bot) CheckRequirements() error {
	if bot.APIToken == "" || bot.Owner <= 0 {
		return errors.New("Either the APIToken, Owner, or both were not found for this bot. Please check the secrets.json or env.conf files")
	}
	return nil
}
