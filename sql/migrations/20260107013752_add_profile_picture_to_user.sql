-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN profile_image VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN profile_image;
-- +goose StatementEnd
