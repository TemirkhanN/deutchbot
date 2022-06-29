package internal

import (
	"DeutchBot/package/cbus"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeutchBot struct {
	commandsHandler cbus.Bus
	responder       *TelegramOutput
}

func NewBot(api *tgbotapi.BotAPI) *DeutchBot {
	responder := NewTgOutputWriter(api)
	commandBus := cbus.NewCommandBus(responder)

	commandBus.RegisterHandler(
		cbus.NewHandlerDefinition(
			NewQuiz(10),
			func(i cbus.Input) bool {
				signal := ResolveSignal(string(i))

				return IsTestRelatedSignal(signal)
			},
		),
	)

	return &DeutchBot{
		commandsHandler: commandBus,
		responder:       responder,
	}
}

func (db *DeutchBot) Consume(message *tgbotapi.Message) {
	db.responder.SwitchChat(message.Chat.ID)

	db.commandsHandler.Handle(cbus.Input(NewRawSignal(message.Chat.ID, message.Text)))

	db.responder.Flush()
}
