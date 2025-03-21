package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	users.Mut.Lock()
	defer users.Mut.Unlock()

	lang := string(data)

	user, ok := ensureUserExists(ctx, b, mes.Message.Chat.ID, lang)
	if !ok {
		return
	}

	user.SendingKey = true
	users.U[user.ChatID] = user

	user.SendMessage(ctx, b, text.SendKey)
}

func onCreateLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	lang := string(data)

	user, ok := ensureUserExists(ctx, b, mes.Message.Chat.ID, lang)
	if !ok {
		return
	}

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
				user.SendMessage(ctx, b, text.CreatingLobbyError)
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
	users.U[user.ChatID] = user
	users.Mut.Unlock()

	kb := inline.New(b).
		Row().
		Button(text.Convert(mes.Message.From.LanguageCode, text.GMafia), []byte(fmt.Sprintf("%d", text.GMafia) + " " + user.Lang), onGameSelect).
		Row().
		Button(text.Convert(mes.Message.From.LanguageCode, text.GBunker), []byte(fmt.Sprintf("%d", text.GBunker) + " " + user.Lang), onGameSelect)

	user.SendReplayMarkup(ctx, b, kb, text.KeyCreatedF, key)

	logger.Log.Info("new lobby is added", zap.String("lobby_key", key))
}

func onGameSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	reformData := strings.Split(string(data), " ")
	gameType, _ := strconv.Atoi(string(reformData[0]))
	lang := reformData[1]

	user, ok := ensureUserExists(ctx, b, mes.Message.Chat.ID, lang)
	if !ok {
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