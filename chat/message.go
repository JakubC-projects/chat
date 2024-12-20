package chat

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id      int         `json:"id"`
	ChatUid uuid.UUID   `json:"chatUid"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`

	SentAt    time.Time `json:"sentAt"`
	ChangedAt time.Time `json:"changedAt"`
	IsDeleted bool      `json:"isDeleted"`
}

type MessageType string
