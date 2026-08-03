-- +goose Up
SET lock_timeout = '4s';

ALTER TABLE refresh_tokens
    ADD COLUMN is_revoked boolean NOT NULL DEFAULT FALSE;

CREATE INDEX idx_refresh_tokens_active_token ON refresh_tokens (token)
WHERE
    is_revoked = FALSE;

-- +goose Down
SET lock_timeout = '4s';

ALTER TABLE refresh_tokens
    DROP COLUMN IF EXISTS is_revoked;

