-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets
(
    key text PRIMARY KEY,
    value text NOT NULL,
    version_id INTEGER DEFAULT 0,
    creation_date text,
    updating_date text,
    deletion_mark INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
