package game

import (
	"context"
	"fmt"

	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

const (
	MINIMUM_MEMBERS_FOR_MAFIA	= 4
	MinimumMembersForBunker		= 4
)

const (
	GMafia = iota
	GBunker
)

func GetGame(lang string, game int) string {
	switch lang {
	case "en":
		return [...]string{"Mafia", "Shelter"}[game]
	case "ru":
		return [...]string{"Мафиа", "Бункер"}[game]
	default:
		return [...]string{"Mafia", "Shelter"}[game]
	}
}

type GameInterface interface {
	startGame(ctx context.Context, b Bot)
}

type Player interface {
	GetInfo() string
}

type User struct {
	ChatID		int64
	Lang		string
	Name		string
	
	LobbyID		int64
	LobbyKey	string
	SendingKey	bool

	player		Player
}

func (u *User) SendMessage(ctx context.Context, b Bot, key int, formats ...any) (*models.Message, error) {
	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:	u.ChatID,
		Text:	u.GetText(key, formats...),
	})
}

func (u *User) SendReplayMarkup(ctx context.Context, b Bot, rm models.ReplyMarkup, key int, formats ...any) (*models.Message, error) {
	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:			u.ChatID,
		Text:			u.GetText(key, formats...),
		ReplyMarkup:	rm,
	})
}

func (u *User) GetText(key int, formats ...any) string {
	switch u.Lang {
	case "en":
		return fmt.Sprintf(text.En[key], formats...)
	case "ru":
		return fmt.Sprintf(text.Ru[key], formats...)
	default:
		return fmt.Sprintf(text.En[key], formats...)
	}
}

type Users		[]User
type UsersRef	[]*User

func (u *Users) sendAll(ctx context.Context, b Bot, key int, a ...any) {
	for _, player := range *u {
		player.SendMessage(ctx, b, key, a...)
	}
}

func (u *Users) findMember(user User) int {
	index := -1

	for i, member := range *u {
		if user.ChatID == member.ChatID {
			index = i
		}
	}

	return index
}

func (u *Users) getMember(chatID int64) (*User, bool) {
	for _, user := range *u {
		if chatID == user.ChatID {
			return &user, true
		}
	}
	return &User{}, false
}

// func (u UsersRef) sendAll(ctx context.Context, b *bot.Bot, msg string, a ...any) {
// 	for _, player := range u {
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID: player.ChatID,
// 			Text:   fmt.Sprintf(msg, a...),
// 		})
// 	}
// }

type Lobby struct {
	LeaderID	int64
	ID			int64

	GameType	int
	IsStarted	bool
	Game		GameInterface

	Members		Users
}

func (l *Lobby) StartGame(ctx context.Context, b *bot.Bot) {
	switch l.GameType {
	case GMafia:
		l.Game = &MafiaGame{
			isStarted:	&l.IsStarted,
			members:	&l.Members,
		}

		for i := range l.Members {
			l.Members[i].player = &MafiaPlayer{isAlive:	true, lang: &l.Members[i].Lang}
		}

	case GBunker:
		l.Game = &BunkerGame{
			isStarted:	&l.IsStarted,
			members:	&l.Members,
		}

		for i := range l.Members {
			l.Members[i].player = &BunkerPlayer{lang: &l.Members[i].Lang}
		}
	}

	l.Game.startGame(ctx, b)
}