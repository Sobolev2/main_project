
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_message_notification
        UNIQUE (type, message_id)
);

CREATE INDEX IF NOT EXISTS idx_notifications_user
ON notifications (user_id, created_at DESC);
