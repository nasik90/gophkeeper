-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets(   
    guid TEXT PRIMARY KEY,
    key BLOB,
    value BLOB,
    binary_value INTEGER DEFAULT 0,
    version_id INTEGER DEFAULT 0,
    creation_date TEXT,
    updating_date TEXT,
    deletion_mark INTEGER DEFAULT 0,
    to_send INTEGER DEFAULT 0,
    comment TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
