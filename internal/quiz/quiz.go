package quiz

import (
	"DeutchBot/internal"
	"DeutchBot/internal/database"
	"DeutchBot/package/cbus"
	"encoding/json"
	"strings"
)

type quiz struct {
	CurrentTask int
	Tasks       []database.Task
}

func (q *quiz) getActiveTask() database.Task {
	return q.Tasks[q.CurrentTask]
}

func (q *quiz) toNextTask() database.Task {
	nextTaskIndex := q.CurrentTask + 1

	if nextTaskIndex > len(q.Tasks)-1 {
		return database.EmptyTask
	}

	q.CurrentTask = nextTaskIndex

	return q.Tasks[q.CurrentTask]
}

func (q quiz) saveQuiz(c *database.Chat) {
	serialized, _ := json.Marshal(q)

	c.SetPayload(string(serialized))
}

func getQuiz(c *database.Chat) *quiz {
	var q quiz
	json.Unmarshal([]byte(c.Payload()), &q)

	return &q
}

func newQuiz(amountOfTasks uint) *quiz {
	q := &quiz{
		CurrentTask: 0,
		Tasks:       []database.Task{},
	}
	iterations := 0
	tasks := make(map[int]database.Task, amountOfTasks)

	for {
		iterations++
		if len(tasks) == int(amountOfTasks) || iterations > 100 {
			break
		}

		task := database.TaskRepository.GetRandom()

		_, duplicate := tasks[task.Id]
		if !duplicate {
			tasks[task.Id] = task
			q.Tasks = append(q.Tasks, task)
		}
	}

	return q
}

type QuizHandler struct {
	amountOfTasks uint
}

func NewQuizHandler(amountOfTasks uint) *QuizHandler {
	return &QuizHandler{amountOfTasks: amountOfTasks}
}

func (qh *QuizHandler) Handle(i cbus.Input, o cbus.Output) {
	s := internal.ResolveSignal(string(i))
	if !CanHandle(s) {
		return
	}

	chat := database.ChatRepository.FindChatById(s.ChatId)
	if chat == nil {
		chat, _ = database.NewChat(s.ChatId, workflowName)
	}

	if chat.State() != stateInProgress {
		qh.start(chat, o)

		return
	}

	qh.applyAnswer(chat, s.Text, o)
}

func (qh *QuizHandler) start(chat *database.Chat, o cbus.Output) {
	chat.SwitchState(stateInProgress)

	chatQuiz := newQuiz(qh.amountOfTasks)
	chatQuiz.saveQuiz(chat)

	currentTask := chatQuiz.getActiveTask()
	o.Write("QuizHandler started.")
	o.Write(currentTask.Question)
}

func (qh *QuizHandler) complete(chat *database.Chat, o cbus.Output) {
	if chat.State() != stateInProgress {
		return
	}

	chat.SwitchState(stateComplete)

	o.Write("Quiz is complete.")
}

func (qh *QuizHandler) applyAnswer(chat *database.Chat, answer string, o cbus.Output) {
	chatQuiz := getQuiz(chat)

	currentTask := chatQuiz.getActiveTask()

	if currentTask.IsCorrectAnswer(answer) {
		o.Write("Correct.")
	} else {
		o.Write("Incorrect.Correct was: " + strings.Join(currentTask.Answers, "; "))
	}

	nextTask := chatQuiz.toNextTask()
	chatQuiz.saveQuiz(chat)
	if database.EmptyTask.Id == nextTask.Id {
		qh.complete(chat, o)

		return
	}

	o.Write(nextTask.Question)
}

func CanHandle(s internal.Signal) bool {
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
	workflowName    = "QuizHandler"
)
