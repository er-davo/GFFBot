package storage

import (
	"time"

	"gffbot/internal/config"
	"gffbot/internal/database"
	"gffbot/internal/logger"
	"gffbot/pkg/semaphore"

	"go.uber.org/zap"
)

func init() {
	logger.Init()
}

type Repository struct {
	db database.DB
	sem *semaphore.Semaphore
}

func NewRepository(db database.DB) *Repository {
    return &Repository{
		db: db,
	    sem: semaphore.New(config.Load().SemaphoreTickets),
    }
}

func (r *Repository) deleteInatciveUsers(days int64) error {
	r.sem.Acquire()
	defer r.sem.Release()
	
	query := `DELETE FROM users WHERE last_activity_game < NOW() INTERVAL '1 day' * $1`
	_, err := r.db.Exec(query, days)
	return err
}

func (r *Repository) CleanUpTask(interval time.Duration, days int64) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		logger.Log.Info("cleaning up task...")
		err := r.deleteInatciveUsers(days)
		if err != nil {
            logger.Log.Error("failed to clean up inactive users", zap.Error(err))
		}
		logger.Log.Info("cleaning up task completed")		
	}
}