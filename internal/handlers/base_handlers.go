package handlers

import (
	"context"
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

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	index, exists := users.FindUser(game.User{ChatID: update.Message.Chat.ID})

	users.Mut.Lock()
	defer users.Mut.Unlock()
	if exists && users.U[index].SendingKey {
		u := users.U[index]

		if u.LobbyID != 0 || u.LobbyKey != "" {
			u.SendMessage(ctx, b, text.AlreadyInLobbyF, u.LobbyKey)
			return
		}

		key := update.Message.Text

		lobbies.Mut.Lock()
		defer lobbies.Mut.Unlock()
		if lob, exists := lobbies.L[key]; exists {

			if lob.IsStarted {
				u.SendMessage(ctx, b, text.LobbyGameIsStarted)
				return
			}

			// Поменять

			memebersList := []string{u.Name}

			for _, memeber := range lob.Members {
				memebersList = append(memebersList, memeber.Name)
			}

			newList := strings.Join(memebersList, "\n")

			// ==//==

			for _, member := range lob.Members {
				member.SendMessage(ctx, b, text.PlayerJoinedLobbyF, u.Name, newList)
			}

			u.SendingKey = false
			u.LobbyKey = key
			u.LobbyID = lob.ID
			users.U[index] = u

			lob.Members = append(lob.Members, u)

			lobbies.L[key] = lob

			u.SendMessage(ctx, b, text.PlayerJoinedLobbyF, u.GetText(text.You), newList)
		} else {
			u.SendMessage(ctx, b, text.LobbyNotExists)
			logger.Log.Error("lobby does not exists", zap.String("lobby_key", u.LobbyKey), zap.Int64("chat_id", u.ChatID))
			return
		}
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.Convert(update.Message.From.LanguageCode, text.UnknownCommand),
		})
		logger.Log.Error("user is not in data", zap.Int64("chat_id", update.Message.Chat.ID))
	}
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button(text.Convert(update.Message.From.LanguageCode, text.Join), []byte("1"), onJoinLobbySelect).
		Row().
		Button(text.Convert(update.Message.From.LanguageCode, text.Create), []byte("2"), onCreateLobbySelect)

	newUser := game.User{
		ChatID: update.Message.Chat.ID,
		Name: bot.EscapeMarkdown(update.Message.From.FirstName) + " " +
			bot.EscapeMarkdown(update.Message.From.LastName),
		Lang: update.Message.From.LanguageCode,
	}

	users.Append(newUser)

	logger.Log.Info("new user added", zap.String("name", newUser.Name), zap.Int64("chat_id", newUser.ChatID))

	newUser.SendReplayMarkup(ctx, b, kb, text.StartCommandF, newUser.Name)
}

func GameStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	currentUser, exists := users.GetUser(update.Message.Chat.ID)
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.Convert(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		logger.Log.Error("user is not in data", zap.Int64("chat_id", update.Message.Chat.ID))
		return
	}

	lobbies.Mut.RLock()
	lob, exists := lobbies.L[currentUser.LobbyKey]
	lobbies.Mut.RUnlock()

	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.Convert(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		logger.Log.Error("lobby does not exists", zap.String("lobby_key", currentUser.LobbyKey), zap.Int64("chat_id", currentUser.ChatID))
		return
	}

	switch lob.GameType {
	case text.GMafia:
		if len(lob.Members) < game.MinimumMembersForMafia {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text.Convert(update.Message.From.LanguageCode, text.AtLeastMembersF, game.MinimumMembersForMafia),
			})
			return
		}

		if len(lob.Members) > game.MaximumMembersForMafia {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text.Convert(update.Message.From.LanguageCode, text.MaximumMembersF, game.MaximumMembersForMafia),
			})
			return
		}
	case text.GBunker:
		if len(lob.Members) < game.MinimumMembersForBunker {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text.Convert(update.Message.From.LanguageCode, text.AtLeastMembersF, game.MinimumMembersForBunker),
			})
			return
		}

		if len(lob.Members) > game.MaximumMembersForBunker {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text.Convert(update.Message.From.LanguageCode, text.MaximumMembersF, game.MaximumMembersForBunker),
			})
			return
		}
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.Convert(update.Message.From.LanguageCode, text.CantStartGame),
		})
		logger.Log.Info("GameType is not selected", zap.String("lobby_key", currentUser.LobbyKey))
		return
	}

	for _, member := range lob.Members {
		member.SendMessage(ctx, b, text.GameStarted)
	}

	lob.IsStarted = true

	lobbies.Mut.Lock()
	lobbies.L[currentUser.LobbyKey] = lob
	lobbies.Mut.Unlock()

	// go lobby.StartGame(ctx, b) ???

	logger.Log.Info("new game started [NEED LOG CHANGES]", zap.Int("game_type", lob.GameType))

	lob.StartGame(ctx, b)
}
