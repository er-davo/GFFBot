package storage

import (
	"database/sql"
	"errors"
)

type User struct {
	ID     int64
	ChatID int64
	Name   string
}

var ErrAlreadyInDatabase = errors.New("user with that chat_id already exists")
var ErrNotFound = errors.New("user with that chat_id or id not found")

func (r *Repository) CreateUser(u User) error {
	r.sem.Acquire()
	defer r.sem.Release()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var userID int64
	userQuery := `INSERT INTO users (chat_id, name) VALUES ($1, $2) ON CONFLICT (chat_id) DO NOTHING RETURNING id`
	err = tx.QueryRow(userQuery, u.ChatID, u.Name).Scan(&userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAlreadyInDatabase
		}
		return err
	}

	statQuery := `INSERT INTO statistic (user_id, played_games, wins, losses, winrate) VALUES ($1, 0, 0, 0, 0.0)`
	_, err = tx.Exec(statQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetUser(chatID int64) (User, error) {
	r.sem.Acquire()
	defer r.sem.Release()

	user := User{ChatID: chatID}
	query := `SELECT id, name FROM users WHERE chat_id = $1`
	row := r.db.QueryRow(query, chatID)
	err := row.Scan(&user.ID, &user.Name)
	return user, err
}

func (r *Repository) UpdateUserActivity(userID int64) error {
	r.sem.Acquire()
	defer r.sem.Release()

	query := `UPDATE users SET last_activity_game = NOW() WHERE id = $1`
	result, err := r.db.Exec(query, userID)
	if err != nil {
        return err
    }

	rowsAffected, err := result.RowsAffected()
	if err != nil {
        return err
    }
	if rowsAffected != 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteUser(u User) error {
	r.sem.Acquire()
	defer r.sem.Release()

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, u.ID)

	return err
}