-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets_update_info (
    data_version TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets_update_info;
-- +goose StatementEnd
