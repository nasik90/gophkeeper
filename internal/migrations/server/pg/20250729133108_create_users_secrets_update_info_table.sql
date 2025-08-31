-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_secrets_update_info (
    user_id int PRIMARY KEY,
    updating_date timestamp NOT NULL,
    update_version int DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_secrets_update_info;
-- +goose StatementEnd
