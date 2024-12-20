package postgres

import (
	"context"
	"testing"

	"github.com/JakubC-projects/chat/chat"
	"github.com/stretchr/testify/assert"
)

var testConnString = "dbname=chat"

func TestListen(t *testing.T) {
	db := NewDb(testConnString)

	listener, err := db.Subscribe(context.Background())
	assert.NoError(t, err)
	defer listener.Close(context.Background())

	eventsReceived := []chat.Event{}
	for range 2 {
		event, err := listener.NextEvent(context.Background())
		assert.NoError(t, err)
		eventsReceived = append(eventsReceived, event)
	}

	assert.NotEmpty(t, eventsReceived)
}

func TestNotify(t *testing.T) {
	db := NewDb(testConnString)
	err := db.sendNotification(context.Background(), "message", "hello1")
	assert.NoError(t, err)

	err = db.sendNotification(context.Background(), "message", "world")
	assert.NoError(t, err)
}

func (d *DB) sendNotification(ctx context.Context, name, msg string) error {
	_, err := d.connPool.Exec(ctx, "SELECT pg_notify($1, $2);", name, msg)
	return err
}
