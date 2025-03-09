package handlers

import (
	"gffbot/internal/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	u := Users{
		U: map[int64]game.User{
			0: {ChatID: 0},
			1: {ChatID: 1},
			2: {ChatID: 2},
		},
	}

	user, ok := u.U[1]

	assert.Equal(t, game.User{ChatID: 1}, user, "Got wrong user index")
	assert.True(t, ok, "Got no user")

	user, ok = u.U[-1]

	assert.Equal(t, game.User{}, user, "Got user index")
	assert.False(t, ok, "Got some user")
}
