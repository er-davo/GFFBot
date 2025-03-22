package database

import (
	"database/sql"
	"fmt"
	
	"gffbot/internal/config"

	_ "github.com/lib/pq"
)

type DB interface {
	Exec(query string, args...interface{}) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
}

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		config.Load().DatabaseHost,
		config.Load().DatabasePort, 
		config.Load().DatabaseUser, 
		config.Load().DatabasePassword, 
		config.Load().DatabaseName,
	)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
        return nil, err
    }
	
	err = db.Ping()
	if err != nil {
        return nil, err
    }

	return db, nil
}