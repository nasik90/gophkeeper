-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets_update_info (
    updating_date text NOT NULL,
    update_version INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets_update_info;
-- +goose StatementEnd
