package quiz

import (
	"DeutchBot/internal/chat"
	"DeutchBot/internal/learn"
	"encoding/json"
)

type quiz struct {
	CurrentTask int
	Tasks       []uint
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
	serialized, _ := json.Marshal(q)

	c.SetPayload(string(serialized))
	chat.ChatRepository.SaveChat(c)
}

func getQuiz(c *chat.Chat) *quiz {
	var q quiz
	json.Unmarshal([]byte(c.Payload), &q)

	return &q
}

func newQuiz(amountOfTasks uint) *quiz {
	q := &quiz{
		CurrentTask: 0,
		Tasks:       learn.TaskRepository.FindRandomTasksIds(amountOfTasks),
	}

	return q
}
