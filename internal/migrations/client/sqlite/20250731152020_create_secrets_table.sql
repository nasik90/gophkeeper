-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets
(   
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    id_on_server INTEGER DEFAULT 0,
    key text PRIMARY KEY,
    value text NOT NULL,
    binary_value INTEGER DEFAULT 0,
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
