package chat

import (
	"DeutschBot/internal"
	"errors"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Workflow string
	State    int
	Payload  string
}

func (c Chat) Id() uint {
	return c.ID
}

func (c Chat) IsInState(state int) bool {
	return c.State == state
}

func (c *Chat) SwitchState(state int) {
	c.State = state
	ChatRepository.SaveChat(c)
}

func (c *Chat) SetPayload(newPayload string) {
	c.Payload = newPayload
	ChatRepository.SaveChat(c)
}

func NewChat(id uint, workflow string) (*Chat, error) {
	if chatExists(id) {
		return nil, errors.New("chat with this id already exists")
	}

	chat := &Chat{
		Model: gorm.Model{
			ID: id,
		},
		Workflow: workflow,
		State:    0,
		Payload:  "",
	}

	ChatRepository.SaveChat(chat)

	return chat, nil
}

func chatExists(id uint) bool {
	if id == 0 {
		return false
	}

	return ChatRepository.FindChatById(id) != nil
}

type chatRepository struct {
	db *gorm.DB
}

func newChatRepository(db *gorm.DB) *chatRepository {
	db.AutoMigrate(&Chat{})

	return &chatRepository{db: db}
}

func (cr *chatRepository) FindChatById(id uint) *Chat {
	var result *Chat
	cr.db.Find(&result, id)

	if result == nil || result.Id() == 0 {
		return nil
	}

	return result
}

func (cr *chatRepository) SaveChat(chat *Chat) {
	if !chatExists(chat.Id()) {
		cr.db.Create(chat)
	} else {
		cr.db.Model(chat).Updates(*chat)
	}
}

var (
	ChatRepository = newChatRepository(internal.Db)
)
