-- +goose Up
-- +goose StatementBegin
CREATE TABLE shorturls (
 id UUID PRIMARY KEY,
 link VARCHAR(255) NOT NULL,
 name VARCHAR(255) NOT NULL,
 created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
 updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shorturls;
-- +goose StatementEnd
