package handlers

import (
	"context"
	"math/rand/v2"
	"sync"
	"time"

	"gffbot/internal/game"
	"gffbot/internal/logger"
	"gffbot/internal/storage"
	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

var users Users
var lobbies Lobbies
var repo *storage.Repository

type Lobbies struct {
	Mut sync.RWMutex
	L   map[string]game.Lobby
}

type Users struct {
	Mut sync.RWMutex
	U   map[int64]game.User
}

func (us *Users) Add(u *game.User) {
	us.Mut.Lock()
	defer us.Mut.Unlock()
	us.U[u.ChatID] = *u
}

func (us *Users) Get(chatID int64) (*game.User, bool) {
	us.Mut.RLock()
    defer us.Mut.RUnlock()
    user, ok := us.U[chatID]
    return &user, ok
}

func createLobbyKey() string {
	key := make([]byte, 4)

	key[0] = text.LettersBytes[rand.IntN(len(text.LettersBytes))]
	key[2] = text.LettersBytes[rand.IntN(len(text.LettersBytes))]
	key[1] = text.DigitsBytes[rand.IntN(len(text.DigitsBytes))]
	key[3] = text.DigitsBytes[rand.IntN(len(text.DigitsBytes))]

	return string(key)
}

func ensureUserExists(ctx context.Context, b *bot.Bot, chatID int64, lang string) (game.User, bool) {
    users.Mut.Lock()
    defer users.Mut.Unlock()
    user, exists := users.U[chatID]
	
    if !exists {
        b.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: chatID,
            Text:   text.Convert(lang, text.CantFindUser),
        })
        logger.Log.Error("user is not in data", zap.Int64("chat_id", chatID))
        return game.User{}, false
    }

	user.Activity = time.Now()
	users.U[user.ChatID] = user

    return user, true
}

func LoadRepository(r *storage.Repository) {
	repo = r
}

func init() {
	lobbies.L = make(map[string]game.Lobby)
	users.U = make(map[int64]game.User)
}
