package handlers

import (
	"gffbot/internal/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	u := Users{
		U: []game.User{
			{ChatID: 0},
			{ChatID: 1},
			{ChatID: 2},
		},
	}

	i, ok := u.FindUser(game.User{ChatID: 1})

	assert.Equal(t, 1, i, "Got wrong user index")
	assert.True(t, ok, "Got no user")

	i, ok = u.FindUser(game.User{ChatID: -1})

	assert.Equal(t, -1, i, "Got user index")
	assert.False(t, ok, "Got some user")
}