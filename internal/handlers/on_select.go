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
	index, exists := users.FindUser(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.ConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	users.Mut.Lock()
	users.U[index].SendingKey = true
	users.Mut.Unlock()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   text.ConvertToLang(mes.Message.From.LanguageCode, text.SendKey),
	})
}

func onGameSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	gameType, _ := strconv.Atoi(string(data))

	user, exists := users.GetUser(mes.Message.Chat.ID)
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.ConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	lobbies.Mut.Lock()
	defer lobbies.Mut.Unlock()
	lobby, exists := lobbies.L[user.LobbyKey]
	if !exists {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		log.Printf("key: %s", user.LobbyKey)
		return
	}

	lobby.GameType = gameType
	lobbies.L[user.LobbyKey] = lobby

	user.SendMessage(ctx, b, text.GameChosenF, text.ConvertToLang(user.Lang, gameType))
}

func onCreateLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := users.FindUser(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.ConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	users.Mut.Lock()
	user := users.U[index]
	users.Mut.Unlock()

	var key string

	lobbies.Mut.Lock()
	defer lobbies.Mut.Unlock()

	if len(lobbies.L) != 0 {
		for i := range 10 {
			key = createLobbyKey()
			if _, exists := lobbies.L[key]; exists {
				break
			} else if i == 9 {
				log.Printf("Somefting went wrong on creating new lobby. Current count of lobbies is: %d", len(lobbies.L))
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   text.ConvertToLang(mes.Message.From.LanguageCode, text.CreatingLobbyError),
				})
				return
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

	lobbies.L[key] = newLobby

	user.LobbyKey = key
	users.Mut.Lock()
	users.U[index] = user
	users.Mut.Unlock()

	kb := inline.New(b).
		Row().
		Button(text.ConvertToLang(mes.Message.From.LanguageCode, text.GMafia), []byte(fmt.Sprintf("%d", text.GMafia)), onGameSelect).
		Row().
		Button(text.ConvertToLang(mes.Message.From.LanguageCode, text.GBunker), []byte(fmt.Sprintf("%d", text.GBunker)), onGameSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        text.ConvertToLang(mes.Message.From.LanguageCode, text.KeyCreatedF, key),
		ReplyMarkup: kb,
	})

	log.Printf("New lobby {ID: %d, key: %s} added", newLobby.ID, key)
}
