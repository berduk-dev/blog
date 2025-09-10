-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    add constraint users_email_key unique (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
