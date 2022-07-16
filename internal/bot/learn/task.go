package learn

import (
	"DeutschBot/internal"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type Example struct {
	Usage   string
	Meaning string
}

type Task struct {
	gorm.Model
	Question string         `gorm:"uniqueIndex:idx_question_answer"`
	Answers  datatypes.JSON `gorm:"uniqueIndex:idx_question_answer"`
	Type     int
	Examples datatypes.JSON
}

func CreateTask(question string, answers []string, examples []Example) *Task {
	task := &Task{
		Question: question,
		Answers:  internal.Serialize(answers),
		Type:     0, // todo
		Examples: internal.Serialize(examples),
	}

	TaskRepository.save(task)

	return task
}

func (t Task) ShowAnswers() []string {
	var answers []string
	internal.Deserialize(t.Answers, &answers)

	return answers
}

func (t Task) ShowExample(number uint) (Example, error) {
	var examples []Example
	internal.Deserialize(t.Examples, &examples)

	if int(number) > len(examples) {
		return Example{}, errors.New(fmt.Sprintf("example %d is out of bounds", number))
	}

	return examples[number-1], nil
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

func (tr *taskRepository) save(task *Task) {
	if task.ID == 0 {
		tr.db.Create(task)
	} else {
		tr.db.Model(task).Updates(*task)
	}
}

var (
	TaskRepository = newTaskRepository(internal.Db)
)
