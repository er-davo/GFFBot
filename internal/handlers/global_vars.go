package handlers

import (
	"math/rand"
	"sync"
	"time"

	"gffbot/internal/game"
	"gffbot/internal/text"
)

var users Users
var lobbies Lobbies

type Lobbies struct {
	Mut sync.RWMutex
	L   map[string]game.Lobby
}

type Users struct {
	Mut sync.RWMutex
	U   []game.User
}

func (us *Users) FindUser(u game.User) (int, bool) {
	us.Mut.RLock()
	defer us.Mut.RUnlock()

	for i, user := range us.U {
		if u.ChatID == user.ChatID {
			return i, true
		}
	}
	return -1, false
}

func (us *Users) GetUser(chatID int64) (game.User, bool) {
	us.Mut.RLock()
	defer us.Mut.RUnlock()

	for _, user := range us.U {
		if chatID == user.ChatID {
			return user, true
		}
	}
	return game.User{}, false
}

func (us *Users) Append(u ...game.User) {
	us.Mut.Lock()
	defer us.Mut.Unlock()
	us.U = append(us.U, u...)
}

func createLobbyKey() string {
	key := make([]byte, 4)

	key[0] = text.LettersBytes[rand.Intn(len(text.LettersBytes))]
	key[2] = text.LettersBytes[rand.Intn(len(text.LettersBytes))]
	key[1] = text.DigitsBytes[rand.Intn(len(text.DigitsBytes))]
	key[3] = text.DigitsBytes[rand.Intn(len(text.DigitsBytes))]

	return string(key)
}

func init() {
	lobbies.L = make(map[string]game.Lobby)
	rand.Seed(time.Now().UnixNano())
}
