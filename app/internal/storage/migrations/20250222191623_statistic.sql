-- +goose Up
-- +goose StatementBegin
CREATE TABLE statistic (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
    played_games INTEGER
    wins INTEGER
    losses INTEGER
    winrate FLOAT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statistic;
-- +goose StatementEnd
