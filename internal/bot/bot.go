package bot

import (
	"DeutchBot/internal/chat"
	"DeutchBot/internal/quiz"
	"DeutchBot/package/cbus"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeutchBot struct {
	commandsHandler cbus.Bus
	responder       *chat.TelegramOutput
}

func NewBot(api *tgbotapi.BotAPI) *DeutchBot {
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

	return &DeutchBot{
		commandsHandler: commandBus,
		responder:       responder,
	}
}

func (db *DeutchBot) Consume(message *tgbotapi.Message) {
	db.responder.SwitchChat(message.Chat.ID)

	db.commandsHandler.Handle(cbus.Input(chat.NewRawSignal(message.Chat.ID, message.Text)))

	db.responder.Flush()
}
