-- +goose Up
-- +goose StatementBegin
ALTER TABLE shorturls
  ADD COLUMN idempotency_key text
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shorturls
  DROP COLUMN IF EXISTS idempotency_key
-- +goose StatementEnd

