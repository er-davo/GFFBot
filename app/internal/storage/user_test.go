package storage

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
        t.Fatalf("an error '%s' on creating sqlmock", err)
    }
	defer db.Close()

	repo := NewRepository(db)

	user := User{
		ChatID: 12345,
        Name:  "Test User",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO users \(chat_id, name\) VALUES \(\$1, \$2\) ON CONFLICT \(chat_id\) DO NOTHING RETURNING id`).
	    WithArgs(user.ChatID, user.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO statistic \(user_id, played_games, wins, losses, winrate\) VALUES \(\$1, 0, 0, 0, 0.0\)`).
	    WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	
	mock.ExpectCommit()

	err = repo.CreateUser(user)

	assert.NoError(t, err, "error creating user")

	assert.NoError(t, mock.ExpectationsWereMet(), "error mock expectations were not met")
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
        t.Fatalf("an error '%s' on creating sqlmock", err)
    }
	defer db.Close()

	repo := NewRepository(db)

	user := User{
		ChatID: 12345,
        Name:  "Test User",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO users \(chat_id, name\) VALUES \(\$1, \$2\) ON CONFLICT \(chat_id\) DO NOTHING RETURNING id`).
	    WithArgs(user.ChatID, user.Name).
		WillReturnError(sql.ErrNoRows)
	
	mock.ExpectRollback()

	err = repo.CreateUser(user)

	assert.ErrorIs(t, err, ErrAlreadyInDatabase, "got an error not equal to ErrAlreadyInDatabase")

	assert.NoError(t, mock.ExpectationsWereMet(), "error mock expectations were not met")
}

// func TestGetUser_Success(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
//         t.Fatalf("an error '%s' on creating sqlmock", err)
//     }
// 	defer db.Close()

// 	repo := NewRepository(db)
// }