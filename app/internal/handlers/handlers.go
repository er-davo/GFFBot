package handlers

import (
	"context"

	"gffbot/internal/game"
	"gffbot/internal/logger"
	"gffbot/internal/storage"
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

	user, ok := ensureUserExists(ctx, b, update.Message.Chat.ID, update.Message.From.LanguageCode)
	if !ok {
		return
	}

	if user.SendingKey {
		KeySubHandler(ctx, b, update, user)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text.Convert(update.Message.From.LanguageCode, text.UnknownCommand),
	})
}

func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text.Convert(update.Message.From.LanguageCode, text.HelpCommand),
	})
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	newUser := game.CreateUser(update)

	users.Add(newUser)

	logger.Log.Info("new user added", zap.String("name", newUser.Name), zap.Int64("chat_id", newUser.ChatID))

	newUser.SendMessage(ctx, b, text.StartCommandF, newUser.Name)
}

func LoginHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, ok := ensureUserExists(ctx, b, update.Message.Chat.ID, update.Message.From.LanguageCode)
	if !ok {
		return
	}

	err := repo.CreateUser(storage.User{
		ChatID: user.ChatID,
		Name:   user.Name,
	})

	if err == storage.ErrAlreadyInDatabase {
		user.SendMessage(ctx, b, text.LoginAlready)
		return
	}

	if err != nil {
		user.SendMessage(ctx, b, text.LoginFailed)
		return
	}

	user.SendMessage(ctx, b, text.LoginSuccess)
}

func StatisticHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, ok := ensureUserExists(ctx, b, update.Message.Chat.ID, update.Message.From.LanguageCode)
	if !ok {
		return
	}

	dataUser, err := repo.GetUser(user.ChatID)
	if err != nil {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		logger.Log.Error("get user failed", zap.Int64("chat_id", update.Message.Chat.ID), zap.Error(err))
		return
	}

	stats, err := repo.GetStatistic(dataUser.ID)
	if err != nil {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		logger.Log.Error("get statistic failed", zap.Int64("chat_id", update.Message.Chat.ID), zap.Error(err))
		return
	}

	user.SendMessage(ctx, b, text.StartCommandF, stats.ToString(user.Lang))
}

func LobbyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button(text.Convert(update.Message.From.LanguageCode, text.Join), []byte(update.Message.From.LanguageCode), onJoinLobbySelect).
		Row().
		Button(text.Convert(update.Message.From.LanguageCode, text.Create), []byte(update.Message.From.LanguageCode), onCreateLobbySelect)

	user, ok := ensureUserExists(ctx, b, update.Message.Chat.ID, update.Message.From.LanguageCode)
	if !ok {
		return
	}

	user.SendReplayMarkup(ctx, b, kb, text.LobbyCommand)
}

func GameStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	currentUser, ok := ensureUserExists(ctx, b, update.Message.Chat.ID, update.Message.From.LanguageCode)
	if !ok {
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

	ok = ensureGame(ctx, b, update, &lob, &currentUser)
	if !ok {
		return
	}

	lob.Members.SendAll(ctx, b, text.GameStarted)

	lob.IsStarted = true

	lobbies.Mut.Lock()
	lobbies.L[currentUser.LobbyKey] = lob
	lobbies.Mut.Unlock()

	logger.Log.Info("new game started [NEED LOG CHANGES]", zap.Int("game_type", lob.GameType))

	lob.StartGame(ctx, b, repo)

	lobbies.Mut.Lock()
	delete(lobbies.L, currentUser.LobbyKey)
	lobbies.Mut.Unlock()
}
