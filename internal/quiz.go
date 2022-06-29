package internal

import (
	"DeutchBot/internal/database"
	"DeutchBot/package/cbus"
	"errors"
	"strings"
)

type Quiz struct {
	currentTask   int
	tasks         []database.Task
	amountOfTasks uint
}

func NewQuiz(amountOfTasks uint) *Quiz {
	return &Quiz{
		currentTask:   0,
		amountOfTasks: amountOfTasks,
		tasks:         []database.Task{},
	}
}

func (q *Quiz) Handle(i cbus.Input, o cbus.Output) {
	s := ResolveSignal(string(i))
	if !IsTestRelatedSignal(s) {
		return
	}

	chat := database.ChatRepository.FindChatById(s.ChatId)
	if chat == nil {
		chat, _ = database.NewChat(s.ChatId, workflowName)
	}

	if chat.State() != stateInProgress {
		q.start(chat, o)

		return
	}

	q.applyAnswer(chat, s.Text, o)
}

func (q *Quiz) start(chat *database.Chat, o cbus.Output) {
	chat.SwitchState(stateInProgress)

	iterations := 0
	tasks := make(map[int]database.Task, q.amountOfTasks)
	q.tasks = []database.Task{}

	for {
		iterations++
		if len(tasks) == int(q.amountOfTasks) || iterations > 100 {
			break
		}

		task := database.TaskRepository.GetRandom()

		_, duplicate := tasks[task.Id]
		if !duplicate {
			tasks[task.Id] = task
			q.tasks = append(q.tasks, task)
		}
	}

	currentTask, _ := q.getActiveTask(chat)
	o.Write("Quiz started.")
	o.Write(currentTask.Question)
}

func (q *Quiz) complete(chat *database.Chat, o cbus.Output) {
	if chat.State() != stateInProgress {
		return
	}

	chat.SwitchState(stateComplete)
	q.currentTask = 0

	o.Write("Quiz is complete.")
}

func (q *Quiz) applyAnswer(chat *database.Chat, answer string, o cbus.Output) {
	currentTask, _ := q.getActiveTask(chat)

	if currentTask.IsCorrectAnswer(answer) {
		o.Write("Correct.")
	} else {
		o.Write("Incorrect.Correct was: " + strings.Join(currentTask.Answers, "; "))
	}

	nextTask := q.toNextTask(chat)
	if database.EmptyTask.Id == nextTask.Id {
		q.complete(chat, o)

		return
	}

	o.Write(nextTask.Question)
}

func (q Quiz) getActiveTask(chat *database.Chat) (database.Task, error) {
	if chat.State() != stateInProgress {
		return database.Task{}, errors.New("quiz is not yet started")
	}

	return q.tasks[q.currentTask], nil
}

func (q *Quiz) toNextTask(chat *database.Chat) database.Task {
	if chat.State() != stateInProgress {
		return database.EmptyTask
	}

	nextTaskIndex := q.currentTask + 1

	if nextTaskIndex > len(q.tasks)-1 {
		return database.EmptyTask
	}

	q.currentTask = nextTaskIndex

	return q.tasks[q.currentTask]
}

func IsTestRelatedSignal(s Signal) bool {
	chat := database.ChatRepository.FindChatById(s.ChatId)
	if chat == nil {
		return s.Text == "/start_quiz"
	}

	if chat.Workflow() == workflowName {
		if chat.State() == stateComplete {
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
	workflowName    = "Quiz"
)
