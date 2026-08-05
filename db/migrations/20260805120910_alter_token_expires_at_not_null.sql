-- +goose Up
ALTER TABLE refresh_tokens
    ALTER COLUMN expires_at SET NOT NULL;

-- +goose Down
ALTER TABLE refresh_tokens
    ALTER COLUMN expires_at DROP NOT NULL;

