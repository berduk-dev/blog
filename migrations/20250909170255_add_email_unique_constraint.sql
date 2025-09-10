-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    add constraint users_email_key unique (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
