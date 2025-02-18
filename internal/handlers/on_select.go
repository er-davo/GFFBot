package handlers

import (
	"context"
	"fmt"
	"strconv"

	"gffbot/internal/game"
	"gffbot/internal/logger"
	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"

	"go.uber.org/zap"
)

func init() {
	logger.Init()
}

func onJoinLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := users.FindUser(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.Convert(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		logger.Log.Error("user is not in data", zap.Int64("chat_id", mes.Message.Chat.ID))
		return
	}

	users.Mut.Lock()
	users.U[index].SendingKey = true
	users.Mut.Unlock()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   text.Convert(mes.Message.From.LanguageCode, text.SendKey),
	})
}

func onGameSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	gameType, _ := strconv.Atoi(string(data))

	user, exists := users.GetUser(mes.Message.Chat.ID)
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.Convert(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		logger.Log.Error("user is not in data", zap.Int64("chat_id", mes.Message.Chat.ID))
		return
	}

	lobbies.Mut.Lock()
	defer lobbies.Mut.Unlock()
	lobby, exists := lobbies.L[user.LobbyKey]
	if !exists {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		logger.Log.Error("lobby does not exists", zap.String("lobby_key", user.LobbyKey), zap.Int64("chat_id", user.ChatID))
		return
	}

	lobby.GameType = gameType
	lobbies.L[user.LobbyKey] = lobby

	user.SendMessage(ctx, b, text.GameChosenF, text.Convert(user.Lang, gameType))
}

func onCreateLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := users.FindUser(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.Convert(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		logger.Log.Error("user is not in data", zap.Int64("chat_id", mes.Message.Chat.ID))
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
				logger.Log.Error("somefting went wrong on creating new lobby.")
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   text.Convert(mes.Message.From.LanguageCode, text.CreatingLobbyError),
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
		Button(text.Convert(mes.Message.From.LanguageCode, text.GMafia), []byte(fmt.Sprintf("%d", text.GMafia)), onGameSelect).
		Row().
		Button(text.Convert(mes.Message.From.LanguageCode, text.GBunker), []byte(fmt.Sprintf("%d", text.GBunker)), onGameSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        text.Convert(mes.Message.From.LanguageCode, text.KeyCreatedF, key),
		ReplyMarkup: kb,
	})

	logger.Log.Info("new lobby is added", zap.String("lobby_key", key))
}
