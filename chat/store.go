package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatStore interface {
	SendMessage(context.Context, Message) (int, error)
	// EditMessage(context.Context, Message) error
	// DeleteMessage(context.Context, Message) error
	GetMessagesBeforeId(ctx context.Context, chatUid uuid.UUID, id int) ([]Message, error)
	GetMessagesAfterDate(ctx context.Context, chatUid uuid.UUID, date time.Time) ([]Message, error)
}
