package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/chat/chat"
	"github.com/jackc/pgx/v5"
)

type EventListener struct {
	conn *pgx.Conn
	db   *DB
}

var _ chat.EventSource = (*EventListener)(nil)

func (d *DB) Subscribe(ctx context.Context) (chat.EventSource, error) {

	conn, err := d.connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire connection for listen: %w", err)
	}

	rawConn := conn.Hijack()

	_, err = rawConn.Exec(ctx, "LISTEN message; LISTEN chat;")
	if err != nil {
		rawConn.Close(ctx)
		return nil, fmt.Errorf("cannot start listening: %w", err)
	}

	return &EventListener{
		conn: rawConn,
		db:   d,
	}, nil
}

func (e *EventListener) NextEvent(ctx context.Context) (chat.Event, error) {
	notification, err := e.conn.WaitForNotification(ctx)
	fmt.Println("New notification")
	if err != nil {
		return chat.Event{}, fmt.Errorf("cannot get next event: %w", err)
	}
	payload := notification.Payload
	if notification.Channel == "message" {
		id, _ := strconv.Atoi(notification.Payload)
		msg, _ := e.db.GetMessage(ctx, id)
		jsonData, _ := json.Marshal(msg)
		payload = string(jsonData)
	}
	return chat.Event{
		Type:    notification.Channel,
		Payload: payload,
	}, nil
}

func (e *EventListener) Close(ctx context.Context) {
	e.conn.Close(ctx)
}
