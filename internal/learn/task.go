package learn

import (
	"DeutchBot/internal"
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type Task struct {
	gorm.Model
	Question string         `gorm:"uniqueIndex:idx_question_answer"`
	Answers  datatypes.JSON `gorm:"uniqueIndex:idx_question_answer"`
}

func (t Task) ShowAnswers() []string {
	var answers []string

	json.Unmarshal(t.Answers, &answers)

	return answers
}

func (t Task) IsCorrectAnswer(answer string) bool {
	for _, applicableAnswer := range t.ShowAnswers() {
		if strings.Title(applicableAnswer) == strings.Title(answer) {
			return true
		}
	}

	return false
}

type taskRepository struct {
	db *gorm.DB
}

func newTaskRepository(db *gorm.DB) *taskRepository {
	db.AutoMigrate(&Task{})

	rand.Seed(time.Now().UnixNano())

	return &taskRepository{db: db}
}

func (tr taskRepository) FindById(id uint) *Task {
	var task *Task

	tr.db.Find(&task, id)

	if task == nil || task.ID == 0 {
		return nil
	}

	return task
}

func (tr taskRepository) FindRandomTasksIds(amount uint) []uint {
	firstTaskId := 1
	var lastTaskId uint

	tr.db.Raw("SELECT MAX(id) FROM tasks").Scan(&lastTaskId)

	tasksIds := make(map[uint]uint)

	for {
		taskId := uint(rand.Intn(int(lastTaskId)) + firstTaskId)
		tasksIds[taskId] = taskId

		if len(tasksIds) == int(amount) {
			break
		}
	}

	ids := make([]uint, 0, len(tasksIds))
	for _, value := range tasksIds {
		ids = append(ids, value)
	}

	return ids
}

func (tr *taskRepository) Save(task *Task) {
	if task.ID == 0 {
		tr.db.Create(task)
	} else {
		tr.db.Model(task).Updates(*task)
	}
}

var (
	TaskRepository = newTaskRepository(internal.Db)
)
