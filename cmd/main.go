package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"gffbot/internal/config"
	"gffbot/internal/handlers"

	"github.com/go-telegram/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	gffbot, err := bot.New(config.Get().TelegramBotApiToken, opts...)
	if err != nil {
		log.Panic(err)
	}

	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/game_start", bot.MatchTypeExact, handlers.GameStartHandler)

	gffbot.Start(ctx)
}