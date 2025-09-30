-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
    guid text NOT NULL PRIMARY KEY,
    key bytea NOT NULL,
    value bytea NOT NULL,
    binary_value BOOLEAN DEFAULT FALSE,
    user_id int NOT NULL,
    version_id BIGINT DEFAULT 1,
    creation_date timestamp NOT NULL,
    updating_date timestamp NOT NULL,
    deletion_mark BOOLEAN DEFAULT FALSE,
    comment text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
