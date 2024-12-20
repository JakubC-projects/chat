package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/JakubC-projects/chat/chat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	db := NewDb(testConnString)

	now := time.Now()

	msg := chat.Message{
		ChatUid:   uuid.New(),
		Type:      "Text",
		Content:   "Hello world",
		SentAt:    now,
		ChangedAt: now,
	}

	msgId, err := db.SendMessage(context.Background(), msg)
	assert.NoError(t, err)
	assert.NotZero(t, msgId)

	messages, err := db.GetMessage(context.Background(), msgId)
	assert.NoError(t, err)
	assert.NotEmpty(t, messages)
}
