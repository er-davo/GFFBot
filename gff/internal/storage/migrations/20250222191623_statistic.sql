-- +goose Up
-- +goose StatementBegin
CREATE TABLE statistic (
    user_id INTEGER references user(id)
    played_games INTEGER
    wins INTEGER
    losses INTEGER
    winrate float
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statistic;
-- +goose StatementEnd
