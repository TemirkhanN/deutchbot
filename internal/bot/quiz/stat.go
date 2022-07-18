package quiz

import (
	"DeutschBot/internal"
	"DeutschBot/internal/bot/learn"
	"gorm.io/gorm"
)

type stat struct {
	gorm.Model
	ChatId         uint
	TaskId         uint
	TookPlace      uint
	CorrectAnswers uint
}

type statService struct {
	repo *gorm.DB
}

func newStatService(db *gorm.DB) *statService {
	db.AutoMigrate(&stat{})

	return &statService{repo: db}
}

func (s statService) LogCorrectAnswer(task learn.Task, inChat uint) {
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

	s.repo.Exec("UPDATE stats SET took_place = took_place + 1, correct_answers = correct_answers + 1 WHERE id= ?", existingStat.ID)
}

func (s statService) LogIncorrectAnswer(task learn.Task, inChat uint) {
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

	s.repo.Exec("UPDATE stats SET took_place = took_place + 1 WHERE id= ?", existingStat.ID)
}

var (
	StatisticService = newStatService(internal.Db)
)
