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

	if !chat.IsInState(stateInProgress) {
		qh.start(chat, o)

		return
	}

	qh.applyAnswer(chat, s.Text, o)
}

func (qh *QuizHandler) start(chat *ch.Chat, o cbus.Output) {
	chat.SwitchState(stateInProgress)

	chatQuiz := newQuiz(qh.amountOfTasks)
	chatQuiz.saveQuiz(chat)

	currentTask := chatQuiz.getActiveTask()
	o.Write("QuizHandler started.")
	o.Write(currentTask.Question)
}

func (qh *QuizHandler) complete(chat *ch.Chat, o cbus.Output) {
	if !chat.IsInState(stateInProgress) {
		return
	}

	chat.SwitchState(stateComplete)

	o.Write("Quiz is complete.")
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

		o.Write("Example")
		o.Write(example.Usage)

		return
	}

	if currentTask.IsCorrectAnswer(answer) {
		o.Write("Correct.")
	} else {
		o.Write("Incorrect.Correct was: " + strings.Join(currentTask.ShowAnswers(), "; "))
	}

	nextTask := chatQuiz.toNextTask()
	chatQuiz.saveQuiz(chat)
	if nextTask == nil {
		qh.complete(chat, o)

		return
	}

	o.Write(nextTask.Question)
}

func CanHandle(s ch.Signal) bool {
	chat := ch.ChatRepository.FindChatById(uint(s.ChatId))
	if chat == nil {
		return s.Text == "/start_quiz"
	}

	if chat.Workflow == workflowName {
		if chat.IsInState(stateComplete) {
			return s.Text == "/start_quiz"
		}

		return true
	}

	return false
}

var (
	stateDraft      = 0
	stateInProgress = 1
	stateComplete   = 2
	workflowName    = "QuizHandler"
)
