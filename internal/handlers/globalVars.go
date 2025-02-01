package handlers

import (
	"gffbot/internal/game"
	"gffbot/internal/text"
	"math/rand"
	"time"
)

var users []game.User
var lobbies map[string]game.Lobby

func findUserInData(u game.User) (int, bool) {
	for i, user := range users {
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

func lastLobbyID() int64 {
	var m int64 = -1

	for _, lobby := range lobbies {
		m = max(m, lobby.ID)
	}

	return m
}

func init() {
	lobbies = make(map[string]game.Lobby)
	rand.Seed(time.Now().UnixNano())
}