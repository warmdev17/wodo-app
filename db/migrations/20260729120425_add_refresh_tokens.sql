-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token varchar(512) NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refresh_token_user_id ON refresh_tokens (user_id);

CREATE INDEX idx_refresh_token_token ON refresh_tokens (token);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;

