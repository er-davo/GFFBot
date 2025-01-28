package game

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKick(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", player: &MafiaPlayer{isAlive: true}},
		User{ChatID: 2, Name: "Player2", player: &MafiaPlayer{isAlive: true}},
	}
	game := &MafiaGame{
		members: &members,
	}
	
	game.kick(members[0])

	if members[0].player.(*MafiaPlayer).isAlive {
		t.Errorf("Expected player to be dead, but they are still alive")
	}
}

func TestMafiaIsDead(t *testing.T) {
	members := Users{
		User{ChatID: 1, Name: "Player1", player: &MafiaPlayer{role: Mafia, isAlive: true}},
		User{ChatID: 2, Name: "Player2", player: &MafiaPlayer{role: Mafia, isAlive: true}},
	}
	game := &MafiaGame{
		members: &members,
		mafias: UsersRef{&members[0], &members[1]},
	}

	if game.mafiaIsDead() {
		t.Errorf("Expected mafia to be alive, got mafia is dead")
	}

	game.mafias[0].player.(*MafiaPlayer).isAlive = false

	if game.mafiaIsDead() {
		t.Errorf("Expected mafia to be alive, got mafia is dead")
	}

	game.mafias[1].player.(*MafiaPlayer).isAlive = false

	if !game.mafiaIsDead() {
		t.Errorf("Expected mafia to be dead, got mafia is alive")
	}
}

func TestSendMessage(t *testing.T) {
	mockBot := new(MockBot)

    // Ожидаем, что SendMessage будет вызван с определенными параметрами
    mockBot.On("SendMessage", mock.Anything, mock.AnythingOfType("*bot.SendMessageParams")).
        Return(&models.Message{Text: "Hello"}, nil)

    // Вызов метода
    message, err := mockBot.SendMessage(context.Background(), &bot.SendMessageParams{Text: "Hello"})


    // Проверка результата
    assert.NoError(t, err)
    assert.NotNil(t, message)
    assert.Equal(t, "Hello", message.Text)

    // Проверка, что ожидания были выполнены
    mockBot.AssertExpectations(t)
}

