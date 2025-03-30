package database

import (
	"database/sql"

	"gffbot/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DB interface {
	Exec(query string, args...interface{}) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
}

func Connect() (*sql.DB, error) {
	psqlURL := config.Load().DatabaseURL
	
	db, err := sql.Open("postgres", psqlURL)
	if err != nil {
        return nil, err
    }
	
	err = db.Ping()
	if err != nil {
        return nil, err
    }

	goose.SetBaseFS(nil)
	err = goose.Up(db, "migrations")
	if err != nil {
        return nil, err
    }

	return db, nil
}