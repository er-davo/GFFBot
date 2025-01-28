package game

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/mock"
)

type MockBot struct {
	mock.Mock
}

func (m *MockBot) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error) {
	args := m.Called(ctx, params)

	return args.Get(0).(*models.Message), args.Error(1)
}