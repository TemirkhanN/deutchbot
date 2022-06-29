package database

import "errors"

type chatRepository struct {
	chats map[int64]*Chat
}

type Chat struct {
	id       int64
	workflow string
	state    int
	payload  string
}

func (c Chat) Id() int64 {
	return c.id
}

func (c Chat) Workflow() string {
	return c.workflow
}

func (c Chat) State() int {
	return c.state
}

func (c *Chat) SwitchState(state int) {
	c.state = state
	ChatRepository.saveChat(c)
}

func (cr *chatRepository) FindChatById(id int64) *Chat {
	return cr.chats[id]
}

func NewChat(id int64, workflow string) (*Chat, error) {
	chat := ChatRepository.FindChatById(id)
	if chat != nil {
		return nil, errors.New("chat with this id already exists")
	}

	chat = &Chat{
		id:       id,
		workflow: workflow,
		state:    0,
	}

	ChatRepository.saveChat(chat)

	return chat, nil
}

func (cr *chatRepository) saveChat(chat *Chat) {
	cr.chats[chat.id] = chat
}

var (
	ChatRepository = &chatRepository{chats: map[int64]*Chat{}}
)
