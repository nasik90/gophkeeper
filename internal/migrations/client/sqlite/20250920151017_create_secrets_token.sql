-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets_token (
    token text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets_token;
-- +goose StatementEnd
