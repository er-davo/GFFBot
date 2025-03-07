-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER SERIAL PRIMARY KEY
    chat_id INTEGER NOT NULL
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
