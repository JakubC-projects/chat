-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: GetMessagesBeforeId :many
SELECT * FROM messages 
WHERE chat_uid = $1 AND id < $2 LIMIT 10;

-- name: GetMessagesAfterDate :many
SELECT * FROM messages 
WHERE chat_uid = $1 AND changed_at > $2;

-- name: CreateMessage :one
INSERT INTO messages (chat_uid, type, content, send_at, changed_at, is_deleted) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: DeleteMessage :exec
UPDATE messages SET is_deleted = true WHERE id = $1;

-- name: NotifyMessageChange :exec
SELECT pg_notify('message', @msg_id);
