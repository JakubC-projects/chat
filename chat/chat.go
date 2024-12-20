package chat

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Uid  uuid.UUID `json:"uid"`
	Name string    `json:"name"`

	MessageCount    int       `json:"messageCount"`
	LastMessage     Message   `json:"lastMessage"`
	LastChangedDate time.Time `json:"lastChangedDate"`
}
