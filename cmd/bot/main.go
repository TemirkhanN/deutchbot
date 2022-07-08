package main

import (
	"DeutchBot/internal/bot"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

// calling without generated tasks will result in error
func main() {
	token := os.Getenv("TG_BOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	deutchBot := bot.NewBot(api)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			deutchBot.Consume(update.Message)
		}
	}
}
