package database

import (
	"embed"
	"gopkg.in/yaml.v3"
	"math/rand"
	"strings"
	"time"
)

type Task struct {
	Id       int
	Question string
	Answers  []string
}

func (t Task) IsCorrectAnswer(answer string) bool {
	for _, applicableAnswer := range t.Answers {
		if strings.Title(applicableAnswer) == strings.Title(answer) {
			return true
		}
	}

	return false
}

type taskRepository struct {
	taskIndex map[int]int
	tasks     []Task
}

func initRepository() *taskRepository {
	dbFile, err := tasksDb.ReadFile("tasks.yaml")
	if err != nil {
		panic(err)
	}

	allTasks := make([]Task, 0)
	_ = yaml.Unmarshal(dbFile, &allTasks)

	repo := &taskRepository{
		taskIndex: make(map[int]int, len(allTasks)),
		tasks:     allTasks,
	}

	for pos, t := range repo.tasks {
		repo.taskIndex[t.Id] = pos
	}

	rand.Seed(time.Now().Unix())

	return repo
}

func (tr taskRepository) GetById(id int) Task {
	return tr.tasks[tr.taskIndex[id]]
}

func (tr taskRepository) GetRandom() Task {
	return tr.tasks[rand.Intn(len(tr.tasks))]
}

var (
	//go:embed tasks.yaml
	tasksDb embed.FS

	EmptyTask      = Task{Id: 0}
	TaskRepository = initRepository()
)
