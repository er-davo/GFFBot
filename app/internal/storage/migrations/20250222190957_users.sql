-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER SERIAL PRIMARY KEY
    chat_id INTEGER NOT NULL UNIQUE
    name VARCHAR(255) NOT NULL
    last_active_game TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
