package quiz

import (
	"DeutschBot/internal"
	"DeutschBot/internal/bot/learn"
	"DeutschBot/internal/chat"
)

type quiz struct {
	CurrentTask int    `json:"CurrentTask"`
	Tasks       []uint `json:"Tasks"`
}

func (q *quiz) getActiveTask() *learn.Task {
	id := q.Tasks[q.CurrentTask]
	if id == 0 {
		return nil
	}

	return learn.TaskRepository.FindById(q.Tasks[q.CurrentTask])
}

func (q *quiz) toNextTask() *learn.Task {
	nextTaskIndex := q.CurrentTask + 1

	if nextTaskIndex > len(q.Tasks)-1 {
		return nil
	}

	q.CurrentTask = nextTaskIndex

	return q.getActiveTask()
}

func (q quiz) saveQuiz(c *chat.Chat) {
	c.SetPayload(string(internal.Serialize(q)))
	chat.ChatRepository.SaveChat(c)
}

func getQuiz(c *chat.Chat) *quiz {
	var q *quiz
	internal.Deserialize([]byte(c.Payload), &q)

	return q
}

func newQuiz(amountOfTasks uint) *quiz {
	q := &quiz{
		CurrentTask: 0,
		Tasks:       learn.TaskRepository.FindRandomTasksIds(amountOfTasks),
	}

	return q
}
