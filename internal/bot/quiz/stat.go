package quiz

import (
	"DeutschBot/internal"
	"DeutschBot/internal/bot/learn"
	"gorm.io/gorm"
)

type stat struct {
	gorm.Model
	ChatId               uint `gorm:"uniqueIndex:idx_task_stat"`
	TaskId               uint `gorm:"uniqueIndex:idx_task_stat"`
	TookPlace            uint
	CorrectAnswers       uint
	LastAttemptSucceeded bool
}

type statService struct {
	repo *gorm.DB
}

func newStatService(db *gorm.DB) *statService {
	db.AutoMigrate(&stat{})

	return &statService{repo: db}
}

func (s statService) logCorrectAnswer(task learn.Task, inChat uint) {
	var existingStat stat
	s.repo.Model(&stat{}).
		Where("task_id = ? AND chat_id = ?", task.ID, inChat).
		Find(&existingStat)

	if existingStat.ID == 0 {
		existingStat = stat{
			ChatId:         inChat,
			TaskId:         task.ID,
			TookPlace:      0,
			CorrectAnswers: 0,
		}

		s.repo.Create(&existingStat)
	}

	s.repo.Exec("UPDATE stats SET took_place = took_place + 1, correct_answers = correct_answers + 1, last_attempt_succeeded = 1 WHERE id= ?", existingStat.ID)
}

func (s statService) logIncorrectAnswer(task learn.Task, inChat uint) {
	var existingStat stat
	s.repo.Model(&stat{}).
		Where("task_id = ? AND chat_id = ?", task.ID, inChat).
		Find(&existingStat)

	if existingStat.ID == 0 {
		existingStat = stat{
			ChatId:         inChat,
			TaskId:         task.ID,
			TookPlace:      0,
			CorrectAnswers: 0,
		}

		s.repo.Create(&existingStat)
	}

	s.repo.Exec("UPDATE stats SET took_place = took_place + 1, last_attempt_succeeded = 0 WHERE id= ?", existingStat.ID)
}

func (s statService) GetMostDifficultTasks(fromChat uint, amount uint) []uint {
	var result []struct {
		TaskId       uint
		WrongAnswers uint
		SuccessRatio float32
	}

	s.repo.Raw(
		"SELECT task_id, (took_place - correct_answers) as wrong_answers, (correct_answers/took_place) as success_ratio "+
			"FROM stats "+
			"WHERE chat_id = ? AND last_attempt_succeeded = 0 "+
			"ORDER BY success_ratio ASC, wrong_answers DESC LIMIT ?",
		fromChat,
		amount,
	).Scan(&result)

	var tasksIds []uint
	for _, taskStat := range result {
		tasksIds = append(tasksIds, taskStat.TaskId)
	}

	return tasksIds
}

var (
	StatisticService = newStatService(internal.Db)
)
