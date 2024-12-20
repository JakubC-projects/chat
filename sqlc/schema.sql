CREATE TABLE IF NOT EXISTS messages
(
    id BIGSERIAL PRIMARY KEY,
    chat_uid uuid NOT NULL,
    type text NOT NULL,
    content text NOT NULL,
    send_at timestamp NOT NULL,
    changed_at timestamp NOT NULL,

    is_deleted boolean NOT NULL
);