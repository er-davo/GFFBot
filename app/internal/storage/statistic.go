package storage

import (
	"fmt"
)

type Statistic struct {
	UserID      int64
	PlayedGames int
	Wins        int
	Losses      int
	Winrate     float64
}

func (s *Statistic) ToString(lang string) string {
	ru := []string{
		"Сыграно игр: %d",
		"Побед: %d",
		"Поражений: %d",
		"Процент побед: %.2f",
	}
	en := []string{
		"Games played: %d",
		"Wins: %d",
        "Losses: %d",
        "Win rate: %.2f",
	}

	if lang == "ru" {
		return fmt.Sprintf(ru[0], s.PlayedGames) + "\n" +
               fmt.Sprintf(ru[1], s.Wins) + "\n" +
               fmt.Sprintf(ru[2], s.Losses) + "\n" +
               fmt.Sprintf(ru[3], s.Winrate*100) + "%"
	} else {
		return fmt.Sprintf(en[0], s.PlayedGames) + "\n" +
               fmt.Sprintf(en[1], s.Wins) + "\n" +
               fmt.Sprintf(en[2], s.Losses) + "\n" +
               fmt.Sprintf(en[3], s.Winrate*100) + "%"
	}
}

func (s *Statistic) Update(win bool) *Statistic {
	s.PlayedGames++
	if win {
		s.Wins++
	} else {
		s.Losses++
	}
	s.Winrate = float64(s.Wins) / float64(s.PlayedGames)

	return s
}

func (r *Repository) GetStatistic(userID int64) (Statistic, error) {
	r.sem.Acquire()
	defer r.sem.Release()

	stat := Statistic{}
	query := `SELECT * FROM statistic WHERE user_id = $1 FOR UPDATE`
	row := r.db.QueryRow(query, userID)
	err := row.Scan(&stat.UserID, &stat.PlayedGames, &stat.Wins, &stat.Losses, &stat.Winrate)
	return stat, err
}

func (r *Repository) UpdateStatistic(stat Statistic) error {
	query := `UPDATE statistic SET played_games = $1, wins = $2, losses = $3, winrate = $4 WHERE user_id = $5`
	_, err := r.db.Exec(query, stat.PlayedGames, stat.Wins, stat.Losses, stat.Winrate, stat.UserID)
	return err
}

func (r *Repository) UpdateUserStatistic(userID int64, win bool) error {
	r.sem.Acquire()
	defer r.sem.Release()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stat, err := r.GetStatistic(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	stat.Update(win)

	err = r.UpdateStatistic(stat)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}