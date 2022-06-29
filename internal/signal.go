package internal

import "encoding/json"

type Signal struct {
	ChatId int64
	Text   string
}

func ResolveSignal(raw string) Signal {
	var result Signal

	json.Unmarshal([]byte(raw), &result)

	return result
}

func NewRawSignal(chatId int64, text string) string {
	data, err := json.Marshal(Signal{
		ChatId: chatId,
		Text:   text,
	})

	if err != nil {
		panic(err)
	}

	return string(data)
}
