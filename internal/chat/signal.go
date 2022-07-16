package chat

import (
	"DeutschBot/internal"
)

type Signal struct {
	ChatId int64
	Text   string
}

func ResolveSignal(raw string) Signal {
	var result Signal
	internal.Deserialize([]byte(raw), &result)

	return result
}

func NewRawSignal(chatId int64, text string) string {
	data := internal.Serialize(Signal{
		ChatId: chatId,
		Text:   text,
	})

	return string(data)
}
