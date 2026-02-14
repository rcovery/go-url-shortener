
-- +goose Up
-- +goose StatementBegin
ALTER TABLE shorturls
  ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '1 day'
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shorturls
  DROP COLUMN IF EXISTS expires_at
-- +goose StatementEnd

