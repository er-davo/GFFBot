package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetStatistic_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
        t.Fatalf("an error '%s' on creating sqlmock", err)
    }
	defer db.Close()

	repo := NewRepository(db)

	expectedStatistic := Statistic{
		UserID:        12345,
        PlayedGames:   10,
        Wins:          8,
        Losses:        2,
        Winrate:       0.8,
	}

	rows := sqlmock.NewRows([]string{"user_id", "played_games", "wins", "losses", "winrate"}).
	    AddRow(expectedStatistic.UserID, expectedStatistic.PlayedGames, expectedStatistic.Wins, expectedStatistic.Losses, expectedStatistic.Winrate)

	mock.ExpectQuery(`SELECT \* FROM statistic WHERE user_id = \$1 FOR UPDATE`).
	    WithArgs(expectedStatistic.UserID).
		WillReturnRows(rows)
	
	stat, err := repo.GetStatistic(expectedStatistic.UserID)

	assert.NoError(t, err, "error getting statistic")

	assert.Equal(t, expectedStatistic, stat, "error getting right statistic")

	assert.NoError(t, mock.ExpectationsWereMet(), "expectations were not met")
}

func TestUpdateUserStatistic_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' on creating sqlmock", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"user_id", "played_games", "wins", "losses", "winrate"}).
		AddRow(1, 5, 3, 2, 0.6)
	mock.ExpectQuery(`SELECT \* FROM statistic WHERE user_id = \$1 FOR UPDATE`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE statistic SET played_games = \$1, wins = \$2, losses = \$3, winrate = \$4 WHERE user_id = \$5`).
		WithArgs(6, 4, 2, 0.6666666666666666, int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.UpdateUserStatistic(1, true)

	assert.NoError(t, err, "error updating user statistic")

	assert.NoError(t, mock.ExpectationsWereMet(), "error mock expectations were not met")
}