package learn

import (
	ch "DeutschBot/internal/chat"
	"DeutschBot/package/cbus"
	"fmt"
	"strings"
)

type WordHandler struct {
}

func (wh WordHandler) Handle(i cbus.Input, o cbus.Output) {
	s := ch.ResolveSignal(string(i))
	if !CanHandle(s) {
		return
	}

	taskId := TaskRepository.FindRandomTasksIds(1)[0]

	task := TaskRepository.FindById(taskId)

	o.Writeln(fmt.Sprintf("Question: %s", task.Question))
	o.Writeln(fmt.Sprintf("Applicable answers: %s.", strings.Join(task.ShowAnswers(), "; ")))
	o.Writeln("Usages:")
	exampleNumber := uint(1)
	for {
		example, err := task.ShowExample(exampleNumber)
		if err != nil {
			break
		}
		o.Writeln(fmt.Sprintf("%d. %s", exampleNumber, example.Usage))
		o.Writeln(example.Meaning)

		exampleNumber++
	}
}

func CanHandle(s ch.Signal) bool {
	chat := ch.ChatRepository.FindChatById(uint(s.ChatId))
	if chat != nil && chat.HasActiveWorkflow() {
		return false
	}

	return s.Text == "/learn_word"
}
