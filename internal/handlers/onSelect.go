package handlers

import (
	"context"
	"fmt"
	"gffbot/internal/game"
	"gffbot/internal/text"
	"log"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func onJoinLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	users[index].SendingKey = true

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.SendKey),
	})
}

func onGameSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	gameType, _ := strconv.Atoi(string(data))

	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	user := users[index]

	lobby, exists := lobbies[user.LobbyKey]
	if !exists {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		log.Printf("key: %s", user.LobbyKey)
		return
	}

	lobby.GameType = gameType
	lobbies[user.LobbyKey] = lobby

	user.SendMessage(ctx, b, text.GameChosenF, game.GetGame(user.Lang, gameType))
}

func onCreateLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	user := users[index]

	var key string

	if len(lobbies) != 0 {
		for i := range 10 {
			key = createLobbyKey()

			if _, exists := lobbies[key]; exists {
				break
			} else if i == 9 {
				log.Printf("Somefting went wrong on creating new lobby. Current count of lobbies is: %d", len(lobbies))
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CreatingLobbyError),
				})
			}
		}
	} else {
		key = createLobbyKey()
	}

	newLobby := game.Lobby{
		LeaderID:  user.ChatID,
		GameType:  text.GameNotSelected,
		IsStarted: false,
		Members:   []game.User{user},
	}

	lobbies[key] = newLobby

	user.LobbyKey = key
	users[index] = user

	kb := inline.New(b).
		Row().
		Button(game.GetGame(mes.Message.From.LanguageCode, game.GMafia), []byte(fmt.Sprintf("%d", game.GMafia)), onGameSelect).
		Row().
		Button(game.GetGame(mes.Message.From.LanguageCode, game.GBunker), []byte(fmt.Sprintf("%d", game.GBunker)), onGameSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        text.GetConvertToLang(mes.Message.From.LanguageCode, text.KeyCreatedF, key),
		ReplyMarkup: kb,
	})

	log.Printf("New lobby {ID: %d, key: %s} added", newLobby.ID, key)
}
