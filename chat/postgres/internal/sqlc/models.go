// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        int64
	ChatUid   uuid.UUID
	Type      string
	Content   string
	SendAt    time.Time
	ChangedAt time.Time
	IsDeleted bool
}
