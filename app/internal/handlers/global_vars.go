package handlers

import (
	"math/rand/v2"
	"sync"

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
	U   map[int64]game.User
}

func (us *Users) Add(u *game.User) {
	us.Mut.Lock()
	defer us.Mut.Unlock()
	us.U[u.ChatID] = *u
}

func createLobbyKey() string {
	key := make([]byte, 4)

	key[0] = text.LettersBytes[rand.IntN(len(text.LettersBytes))]
	key[2] = text.LettersBytes[rand.IntN(len(text.LettersBytes))]
	key[1] = text.DigitsBytes[rand.IntN(len(text.DigitsBytes))]
	key[3] = text.DigitsBytes[rand.IntN(len(text.DigitsBytes))]

	return string(key)
}

func init() {
	lobbies.L = make(map[string]game.Lobby)
	users.U = make(map[int64]game.User)
}
