package main

import (
	"context"
	"os"
	"os/signal"

	"gffbot/internal/config"
	"gffbot/internal/handlers"
	"gffbot/internal/logger"

	"github.com/go-telegram/bot"
)

func main() {
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
	
	logger.Init()
	defer logger.Log.Sync()

	gffbot.Start(ctx)
}
