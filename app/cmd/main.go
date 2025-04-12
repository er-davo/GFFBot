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

	logger.Log.Info("connecting to database...")
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logger.Log.Info("connected succesfull")

	repo := storage.NewRepository(db)
	go repo.CleanUpTask(
		24*time.Hour*time.Duration(config.Load().CheckTimeDatabase),
		config.Load().InactiveDaysDuration,
	)
	handlers.LoadRepository(repo)

	go handlers.UserActivityCleanUp(
		time.Duration(config.Load().CheckTimeMemory),
		time.Minute*time.Duration(config.Load().InactiveMinutsDuration),
	)

	go handlers.LobbyActivityCleanUp(
		time.Duration(config.Load().CheckTimeMemory),
		time.Hour*time.Duration(config.Load().InactiveHoursDuration),
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
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.HelpHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/login", bot.MatchTypeExact, handlers.LoginHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/statistic", bot.MatchTypeExact, handlers.StatisticHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/lobby", bot.MatchTypePrefix, handlers.LobbyHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/game_start", bot.MatchTypeExact, handlers.GameStartHandler)

	logger.Log.Info("starting bot...")
	gffbot.Start(ctx)
}
