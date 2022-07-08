package chat

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramOutput struct {
	chatId int64
	api    *tgbotapi.BotAPI
	buffer bytes.Buffer
}

func NewTgOutputWriter(api *tgbotapi.BotAPI) *TelegramOutput {
	return &TelegramOutput{api: api}
}

func (to *TelegramOutput) SwitchChat(chatId int64) {
	to.chatId = chatId
}

func (to *TelegramOutput) Write(text string) {
	to.buffer.WriteString(text + "\n")
}

func (to *TelegramOutput) Flush() {
	msg := tgbotapi.NewMessage(to.chatId, to.buffer.String())
	to.buffer.Reset()

	to.api.Send(msg)
}
