package bot

import (
	"DeutschBot/internal/bot/learn"
	"DeutschBot/internal/bot/quiz"
	"DeutschBot/internal/chat"
	"DeutschBot/package/cbus"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeutschBot struct {
	commandsHandler cbus.Bus
	responder       *chat.TelegramOutput
}

func NewBot(api *tgbotapi.BotAPI) *DeutschBot {
	responder := chat.NewTgOutputWriter(api)
	commandBus := cbus.NewCommandBus(responder)

	commandBus.RegisterHandler(
		cbus.NewHandlerDefinition(
			quiz.NewQuizHandler(10),
			func(i cbus.Input) bool {
				signal := chat.ResolveSignal(string(i))

				return quiz.CanHandle(signal)
			},
		),
	)

	commandBus.RegisterHandler(
		cbus.NewHandlerDefinition(
			learn.WordHandler{},
			func(i cbus.Input) bool {
				signal := chat.ResolveSignal(string(i))

				return learn.CanHandle(signal)
			},
		),
	)

	return &DeutschBot{
		commandsHandler: commandBus,
		responder:       responder,
	}
}

func (db *DeutschBot) Consume(message *tgbotapi.Message) {
	db.responder.SwitchChat(message.Chat.ID)

	db.commandsHandler.Handle(cbus.Input(chat.NewRawSignal(message.Chat.ID, message.Text)))

	db.responder.Flush()
}
