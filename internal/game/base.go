package game

import (
	"context"
	"fmt"
	"sync"

	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

const (
	MinimumMembersForMafia  = 3
	MaximumMembersForMafia  = 20
	MinimumMembersForBunker = 3
	MaximumMembersForBunker = 22
)

type GameStarter interface {
	StartGame(ctx context.Context, b Bot)
}

type Player interface {
	Info() string
}

type User struct {
	ChatID int64
	Lang   string
	Name   string

	LobbyID    int64
	LobbyKey   string
	SendingKey bool

	Player Player
}

func (u *User) SendMessage(ctx context.Context, b Bot, key int, formats ...any) (*models.Message, error) {
	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.ChatID,
		Text:   u.GetText(key, formats...),
	})
}

func (u *User) SendMessageSync(ctx context.Context, b Bot, key int, wg *sync.WaitGroup, formats ...any) (*models.Message, error) {
	defer wg.Done()
	return u.SendMessage(ctx, b, key, formats...)
}

func (u *User) SendReplayMarkup(ctx context.Context, b Bot, rm models.ReplyMarkup, key int, formats ...any) (*models.Message, error) {
	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      u.ChatID,
		Text:        u.GetText(key, formats...),
		ReplyMarkup: rm,
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

type Users []User
type UsersRef []*User

func (u *Users) SendAll(ctx context.Context, b Bot, key int, a ...any) {
	wg := sync.WaitGroup{}
	wg.Add(len(*u))

	for _, player := range *u {
		go func() {
			defer wg.Done()
			player.SendMessage(ctx, b, key, a...)
		}()
	}

	wg.Wait()
}

func (u *Users) FindMember(user User) int {
	index := -1

	for i, member := range *u {
		if user.ChatID == member.ChatID {
			index = i
		}
	}

	return index
}

func (u *Users) GetMember(chatID int64) (*User, bool) {
	for _, user := range *u {
		if chatID == user.ChatID {
			return &user, true
		}
	}
	return &User{}, false
}

type Lobby struct {
	LeaderID int64
	ID       int64

	GameType  int
	IsStarted bool
	Game      GameStarter

	Members Users
}

func (l *Lobby) StartGame(ctx context.Context, b *bot.Bot) {
	switch l.GameType {
	case text.GMafia:

		l.Game = &MafiaGame{
			IsStarted: &l.IsStarted,
			Members:   &l.Members,
		}

		for i := range l.Members {
			l.Members[i].Player = &MafiaPlayer{Lang: &l.Members[i].Lang}
		}

	case text.GBunker:
		l.Game = &BunkerGame{
			IsStarted: &l.IsStarted,
			Members:   &l.Members,
		}

		for i := range l.Members {
			l.Members[i].Player = &BunkerPlayer{Lang: &l.Members[i].Lang}
		}
	}

	l.Game.StartGame(ctx, b)
}
