package cmd

import (
	"DeutchBot/internal"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	deutchBot := internal.NewBot(bot)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			deutchBot.Consume(update.Message)
		}
	}
}
