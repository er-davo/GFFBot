package handlers

import (
	"gffbot/internal/game"
	"gffbot/internal/text"
	"math/rand"
	"sync"
	"time"
)

var users Users
var lobbies Lobbies

type Lobbies struct {
	Mut sync.Mutex
	L   map[string]game.Lobby
}

type Users struct {
	Mut sync.Mutex
	U   []game.User
}

func (us *Users) findUserInData(u game.User) (int, bool) {
	us.Mut.Lock()
	defer us.Mut.Unlock()

	for i, user := range us.U {
		if u.ChatID == user.ChatID {
			return i, true
		}
	}
	return -1, false
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
