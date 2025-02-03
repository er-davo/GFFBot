package game

import (
	"context"
	"testing"

	"gffbot/internal/botmock"
	"gffbot/internal/text"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFillRoles(t *testing.T) {
	mockBot := new(botmock.MockBot)
	lang := ""

	us := Users{
		User{Player: &MafiaPlayer{Lang: &lang}},
		User{Player: &MafiaPlayer{Lang: &lang}},
		User{Player: &MafiaPlayer{Lang: &lang}},
	}

	mg := MafiaGame{
		Members: &us,
	}

	mockBot.On("SendMessage", mock.Anything, mock.AnythingOfType("*bot.SendMessageParams")).
		Return(&models.Message{Text: ""}, nil)

	mg.fillRoles(context.Background(), mockBot)

	assert.Equal(t, mg.mafias[0].Player.(*MafiaPlayer).role, text.Mafia, "Mafia in not filled")
	assert.Equal(t, mg.detectives[0].Player.(*MafiaPlayer).role, text.Detective, "Detective in not filled")
	assert.Equal(t, mg.doctors[0].Player.(*MafiaPlayer).role, text.Doctor, "Doctor in not filled")
}

func TestKick(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", Player: &MafiaPlayer{isAlive: true}},
		User{ChatID: 2, Name: "Player2", Player: &MafiaPlayer{isAlive: true}},
	}
	game := &MafiaGame{
		Members: &members,
	}

	game.kick(members[0])

	assert.False(t, members[0].Player.(*MafiaPlayer).isAlive,
		"Expected player to be dead, but they are still alive")
}

func TestMafiaIsDead(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", Player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
		User{ChatID: 2, Name: "Player2", Player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
	}
	game := &MafiaGame{
		Members: &members,
		mafias:  UsersRef{&members[0], &members[1]},
	}

	assert.False(t, game.mafiaIsDead(),
		"Expected mafia to be alive, got mafia is dead")

	game.mafias[0].Player.(*MafiaPlayer).isAlive = false

	assert.False(t, game.mafiaIsDead(),
		"Expected mafia to be alive, got mafia is dead")

	game.mafias[1].Player.(*MafiaPlayer).isAlive = false

	assert.True(t, game.mafiaIsDead(),
		"Expected mafia to be dead, got mafia is alive")
}

func TestGetTwoMaxVotes(t *testing.T) {
	users := Users{
		User{Player: &MafiaPlayer{votes: 3}},
		User{Player: &MafiaPlayer{votes: 2}},
		User{Player: &MafiaPlayer{votes: 3}},
	}

	game := MafiaGame{Members: &users}

	maxFirst, maxSecond := game.getTwoMaxVotes()

	assert.Equal(t, users[0].Player, maxFirst.Player, "First max vote player should be users[0].")
	assert.Equal(t, users[2].Player, maxSecond.Player, "Second max vote player should be users[2].")

	users[1].Player.(*MafiaPlayer).votes = 5
	users[0].Player.(*MafiaPlayer).votes = 0

	maxFirst, maxSecond = game.getTwoMaxVotes()

	assert.Equal(t, users[1].Player, maxFirst.Player, "First max vote player should be users[1].")
	assert.Equal(t, users[2].Player, maxSecond.Player, "Second max vote player should be users[2].")
}

func TestCiviliansIsDead(t *testing.T) {
	users := Users{
		User{Player: &MafiaPlayer{role: text.Civilian, isAlive: false}},
		User{Player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
		User{Player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
	}

	game := MafiaGame{Members: &users}

	assert.True(t, game.civiliansIsDead(), "Civilians should be dead.")
}
