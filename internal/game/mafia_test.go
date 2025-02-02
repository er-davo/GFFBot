package game

import (
	"context"
	"testing"

	"gffbot/internal/text"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFillRoles(t *testing.T) {
	mockBot := new(MockBot)
	lang := ""

	us := Users{
		User{player: &MafiaPlayer{lang: &lang}},
		User{player: &MafiaPlayer{lang: &lang}},
		User{player: &MafiaPlayer{lang: &lang}},
	}

	mg := MafiaGame{
		members: &us,
	}

	mockBot.On("SendMessage", mock.Anything, mock.AnythingOfType("*bot.SendMessageParams")).
		Return(&models.Message{Text: ""}, nil)

	mg.fillRoles(context.Background(), mockBot)

	assert.Equal(t, mg.mafias[0].player.(*MafiaPlayer).role, text.Mafia, "Mafia in not filled")
	assert.Equal(t, mg.detectives[0].player.(*MafiaPlayer).role, text.Detective, "Detective in not filled")
	assert.Equal(t, mg.doctors[0].player.(*MafiaPlayer).role, text.Doctor, "Doctor in not filled")
}

func TestKick(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", player: &MafiaPlayer{isAlive: true}},
		User{ChatID: 2, Name: "Player2", player: &MafiaPlayer{isAlive: true}},
	}
	game := &MafiaGame{
		members: &members,
	}
	
	game.kick(members[0])

	assert.False(t, members[0].player.(*MafiaPlayer).isAlive,
		"Expected player to be dead, but they are still alive")
}

func TestMafiaIsDead(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
		User{ChatID: 2, Name: "Player2", player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
	}
	game := &MafiaGame{
		members: &members,
		mafias: UsersRef{&members[0], &members[1]},
	}

	assert.False(t, game.mafiaIsDead(),
		"Expected mafia to be alive, got mafia is dead")

	game.mafias[0].player.(*MafiaPlayer).isAlive = false

	assert.False(t, game.mafiaIsDead(),
		"Expected mafia to be alive, got mafia is dead")

	game.mafias[1].player.(*MafiaPlayer).isAlive = false

	assert.True(t, game.mafiaIsDead(),
		"Expected mafia to be dead, got mafia is alive")
}

func TestGetTwoMaxVotes(t *testing.T) {
	users := Users{
		User{player: &MafiaPlayer{votes: 3}},
		User{player: &MafiaPlayer{votes: 2}},
		User{player: &MafiaPlayer{votes: 3}},
	}

	game := MafiaGame{members: &users}

	maxFirst, maxSecond := game.getTwoMaxVotes()

	assert.Equal(t, users[0].player, maxFirst.player, "First max vote player should be users[0].")
	assert.Equal(t, users[2].player, maxSecond.player, "Second max vote player should be users[2].")

	users[1].player.(*MafiaPlayer).votes = 5
	users[0].player.(*MafiaPlayer).votes = 0

	maxFirst, maxSecond = game.getTwoMaxVotes()

	assert.Equal(t, users[1].player, maxFirst.player, "First max vote player should be users[1].")
	assert.Equal(t, users[2].player, maxSecond.player, "Second max vote player should be users[2].")
}

func TestCiviliansIsDead(t *testing.T) {
	users := Users{
		User{player: &MafiaPlayer{role: text.Civilian, isAlive: false}},
		User{player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
		User{player: &MafiaPlayer{role: text.Mafia, isAlive: true}},
	}

	game := MafiaGame{members: &users}

	assert.True(t, game.civiliansIsDead(), "Civilians should be dead.")
}