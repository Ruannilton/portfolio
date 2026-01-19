-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN github_access_token VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN github_access_token;
-- +goose StatementEnd
