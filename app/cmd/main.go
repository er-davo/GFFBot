package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"gffbot/internal/config"
	"gffbot/internal/database"
	"gffbot/internal/handlers"
	"gffbot/internal/logger"
	"gffbot/internal/storage"

	"github.com/go-telegram/bot"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	db, err := database.Connect()
	if err != nil {
        panic(err)
    }
	defer db.Close()

	repo := storage.NewRepository(db)
	go repo.CleanUpTask(
		24 * time.Hour * time.Duration(config.Load().CheckTime), 
		config.Load().InactiveDaysDuration,
	)


	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	
	gffbot, err := bot.New(config.Load().TelegramBotApiToken, opts...)
	if err != nil {
		panic(err)
	}
	
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/game_start", bot.MatchTypeExact, handlers.GameStartHandler)
	

	gffbot.Start(ctx)
}
