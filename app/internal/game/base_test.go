package game

import (
	"context"
	"testing"

	"gffbot/internal/mocks"
	"gffbot/internal/text"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetText(t *testing.T) {
	u := User{Lang: "en"}
	actual := u.GetText(text.Yes)
	
	assert.Equalf(t, "Yes", actual, "wanted: Yes\ngot: %s", actual)

	u.Lang = "ru"
	actual = u.GetText(text.Yes)
	
	assert.Equalf(t, "Да", actual, "wanted: Да\ngot: %s", actual)

	u.Lang = "something else"
	actual = u.GetText(text.Yes)
	
	assert.Equalf(t, "Yes", actual, "wanted: Yes\ngot: %s", actual)
}

func TestUser_SendMessage(t *testing.T) {
	mockBot := new(mocks.MockBot)
	u := User{}

	mockBot.On("SendMessage", mock.Anything, mock.AnythingOfType("*bot.SendMessageParams")).
		Once().Return(&models.Message{Text: ""}, nil)
	
	message, err := u.SendMessage(context.Background(), mockBot, text.Default)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "", message.Text)

    mockBot.AssertExpectations(t)
}

func TestSendAll(t *testing.T) {
	mockBot := new(mocks.MockBot)
	us := Users{
		User{},
		User{},
		User{},
		User{},
		User{},
	}

	mockBot.On("SendMessage", mock.Anything, mock.AnythingOfType("*bot.SendMessageParams")).
		Times(5).Return(&models.Message{Text: ""}, nil)
	
	us.SendAll(context.Background(), mockBot, text.Default)

	mockBot.AssertExpectations(t)
}