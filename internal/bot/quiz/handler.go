package quiz

import (
	ch "DeutschBot/internal/chat"
	"DeutschBot/package/cbus"
	"log"
	"strings"
)

type QuizHandler struct {
	amountOfTasks uint
}

func NewQuizHandler(amountOfTasks uint) *QuizHandler {
	return &QuizHandler{amountOfTasks: amountOfTasks}
}

func (qh *QuizHandler) Handle(i cbus.Input, o cbus.Output) {
	s := ch.ResolveSignal(string(i))
	if !CanHandle(s) {
		return
	}

	chat := ch.ChatRepository.FindChatById(uint(s.ChatId))
	if chat == nil {
		chat, _ = ch.NewChat(uint(s.ChatId), workflowName)
	}

	if !chat.HasActiveWorkflow() {
		qh.start(chat, o)

		return
	}

	qh.applyAnswer(chat, s.Text, o)
}

func (qh *QuizHandler) start(chat *ch.Chat, o cbus.Output) {
	chat.StartWorkflow()
	chatQuiz := newQuiz(chat, qh.amountOfTasks)

	currentTask := chatQuiz.getActiveTask()
	o.Writeln("QuizHandler started.")
	o.Writeln(currentTask.Question)
}

func (qh *QuizHandler) applyAnswer(chat *ch.Chat, answer string, o cbus.Output) {
	chatQuiz := getQuiz(chat)

	currentTask := chatQuiz.getActiveTask()

	if answer == "/example" {
		example, err := currentTask.ShowExample(1)

		if err != nil {
			log.Print(err)

			return
		}

		o.Writeln("Example")
		o.Writeln(example.Usage)

		return
	}

	if currentTask.IsCorrectAnswer(answer) {
		o.Writeln("Correct.")
		StatisticService.logCorrectAnswer(*currentTask, chat.ID)
	} else {
		o.Writeln("Incorrect.Correct was: " + strings.Join(currentTask.ShowAnswers(), "; "))
		StatisticService.logIncorrectAnswer(*currentTask, chat.ID)
	}

	nextTask := chatQuiz.toNextTask()
	chatQuiz.saveQuiz(chat)
	if nextTask == nil {
		qh.complete(chat, o)

		return
	}

	o.Writeln(nextTask.Question)
}

func (qh *QuizHandler) complete(chat *ch.Chat, o cbus.Output) {
	if !chat.HasActiveWorkflow() {
		return
	}

	chat.CompleteWorkflow()

	o.Writeln("Quiz is complete.")
}

func CanHandle(s ch.Signal) bool {
	chat := ch.ChatRepository.FindChatById(uint(s.ChatId))
	if chat == nil {
		return s.Text == "/start_quiz"
	}

	if chat.Workflow == workflowName {
		if !chat.HasActiveWorkflow() {
			return s.Text == "/start_quiz"
		}

		return true
	}

	return false
}

var (
	workflowName = "QuizHandler"
)
