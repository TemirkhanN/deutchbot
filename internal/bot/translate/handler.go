package translate

import (
	"DeutschBot/internal"
	ch "DeutschBot/internal/chat"
	"DeutschBot/package/cbus"
	"DeutschBot/package/translator"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type TranslationHandler struct {
}

func (th TranslationHandler) Handle(i cbus.Input, o cbus.Output) {
	s := ch.ResolveSignal(string(i))
	if !CanHandle(s) {
		return
	}

	word := parseWord(s.Text)

	res, err := internal.Translator.Translate(word, translator.DE, translator.EN)
	if err != nil {
		log.Println(err)
	}

	var meanings []string
	for _, meaning := range res.Meanings() {
		meanings = append(meanings, meaning.Word())
	}

	o.Writeln(fmt.Sprintf("Translation: %s", strings.Join(meanings, "; ")))
	o.Writeln("Usages:")

	for i, example := range res.Examples() {
		o.Writeln(fmt.Sprintf("%d. %s", i+1, example.Usage()))
		o.Writeln(example.Meaning())
	}
}

func parseWord(commandInput string) string {
	pattern := regexp.MustCompile("/translate ([\\p{L}a-zA-Z ]+)")
	res := pattern.FindStringSubmatch(commandInput)

	if len(res) == 0 {
		return ""
	}

	return res[1]
}

func CanHandle(s ch.Signal) bool {
	chat := ch.ChatRepository.FindChatById(uint(s.ChatId))
	if chat != nil && chat.HasActiveWorkflow() {
		return false
	}

	return parseWord(s.Text) != ""
}
