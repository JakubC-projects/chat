package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/JakubC-projects/chat/chat"
	"github.com/JakubC-projects/chat/chat/postgres/internal/sqlc"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (d *DB) SendMessage(ctx context.Context, msg chat.Message) (int, error) {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txq := d.queries.WithTx(tx)

	msgId, err := txq.CreateMessage(ctx, sqlc.CreateMessageParams{
		ChatUid:   msg.ChatUid,
		Type:      string(msg.Type),
		Content:   msg.Content,
		SendAt:    msg.SentAt,
		ChangedAt: msg.ChangedAt,
		IsDeleted: msg.IsDeleted,
	})
	if err != nil {
		return 0, fmt.Errorf("cannot notify create message: %w", err)
	}

	err = txq.NotifyMessageChange(ctx, strconv.Itoa(int(msgId)))
	if err != nil {
		return 0, fmt.Errorf("cannot notify about message change: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("cannot commit message sent transaction: %w", err)
	}

	return int(msgId), nil
}

func (d *DB) DeleteMessage(ctx context.Context, id int) error {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txq := d.queries.WithTx(tx)

	err = txq.DeleteMessage(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("cannot notify delete message: %w", err)
	}

	err = txq.NotifyMessageChange(ctx, strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("cannot notify about message change: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit message sent transaction: %w", err)
	}

	return nil
}

func (d *DB) GetMessagesBeforeId(ctx context.Context, chatUid uuid.UUID, id int) ([]chat.Message, error) {
	msgs, err := d.queries.GetMessagesBeforeId(ctx, sqlc.GetMessagesBeforeIdParams{
		ChatUid: chatUid,
		ID:      int64(id),
	})
	return lo.Map(msgs, mapSqlcMessage), err
}

func (d *DB) GetMessagesAfterDate(ctx context.Context, chatUid uuid.UUID, date time.Time) ([]chat.Message, error) {
	msgs, err := d.queries.GetMessagesAfterDate(ctx, sqlc.GetMessagesAfterDateParams{
		ChatUid:   chatUid,
		ChangedAt: date,
	})
	return lo.Map(msgs, mapSqlcMessage), err
}

func (d *DB) GetMessage(ctx context.Context, id int) (chat.Message, error) {
	msg, err := d.queries.GetMessage(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return chat.Message{}, chat.ErrNotFound
	}
	return mapSqlcMessage(msg, 0), err
}

func mapSqlcMessage(sqlcMsg sqlc.Message, _ int) chat.Message {
	return chat.Message{
		Id:        int(sqlcMsg.ID),
		ChatUid:   sqlcMsg.ChatUid,
		Type:      chat.MessageType(sqlcMsg.Type),
		Content:   sqlcMsg.Content,
		SentAt:    sqlcMsg.SendAt,
		ChangedAt: sqlcMsg.ChangedAt,
		IsDeleted: sqlcMsg.IsDeleted,
	}
}
