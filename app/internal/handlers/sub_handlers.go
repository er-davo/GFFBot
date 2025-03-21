package handlers

import (
	"context"
	"gffbot/internal/game"
	"gffbot/internal/logger"
	"gffbot/internal/text"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func KeySubHandler(ctx context.Context, b *bot.Bot, update *models.Update, user game.User) {
	if user.LobbyKey != "" {
		user.SendMessage(ctx, b, text.AlreadyInLobbyF, user.LobbyKey)
		return
	}

	key := update.Message.Text

	lobbies.Mut.Lock()
	defer lobbies.Mut.Unlock()
	if lob, exists := lobbies.L[key]; exists {
		if lob.IsStarted {
			user.SendMessage(ctx, b, text.LobbyGameIsStarted)
			return
		}

		// Поменять

		memebersList := []string{user.Name}

		for _, memeber := range lob.Members {
			memebersList = append(memebersList, memeber.Name)
		}

		newList := strings.Join(memebersList, "\n")

		// ==//==

		lob.Members.SendAll(ctx, b, text.PlayerJoinedLobbyF, user.Name, newList)
		
		user.SendingKey = false
		user.LobbyKey = key
		user.LobbyID = lob.ID

		users.Mut.Lock()
		users.U[update.Message.Chat.ID] = user
		users.Mut.Unlock()

		lob.Members = append(lob.Members, user)

		lobbies.L[key] = lob

		user.SendMessage(ctx, b, text.PlayerJoinedLobbyF, user.GetText(text.You), newList)
	} else {
		user.SendMessage(ctx, b, text.LobbyNotExists)
		logger.Log.Error("lobby does not exists", zap.String("lobby_key", user.LobbyKey), zap.Int64("chat_id", user.ChatID))
	}
}
