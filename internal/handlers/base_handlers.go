package handlers

import (
	"context"
	"gffbot/internal/game"
	"gffbot/internal/text"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	index, exists := users.FindUser(game.User{ChatID: update.Message.Chat.ID})

	users.Mut.Lock()
	defer users.Mut.Unlock()
	if exists && users.U[index].SendingKey {
		u := users.U[index]

		if u.LobbyID != 0 || u.LobbyKey != "" {
			u.SendMessage(ctx, b, text.AlreadyInLobbyF, u.LobbyKey)
			return
		}

		key := update.Message.Text

		lobbies.Mut.Lock()
		defer lobbies.Mut.Unlock()
		if lob, exists := lobbies.L[key]; exists {

			if lob.IsStarted {
				u.SendMessage(ctx, b, text.LobbyGameIsStarted)
				return
			}

			// Поменять

			memebersList := []string{u.Name}

			for _, memeber := range lob.Members {
				memebersList = append(memebersList, memeber.Name)
			}

			newList := strings.Join(memebersList, "\n")

			// ==//==

			for _, member := range lob.Members {
				member.SendMessage(ctx, b, text.PlayerJoinedLobbyF, u.Name, newList)
			}

			u.SendingKey = false
			u.LobbyKey = key
			u.LobbyID = lob.ID
			users.U[index] = u

			lob.Members = append(lob.Members, u)

			lobbies.L[key] = lob

			u.SendMessage(ctx, b, text.PlayerJoinedLobbyF, u.GetText(text.You), newList)
		} else {
			u.SendMessage(ctx, b, text.LobbyNotExists)
			return
		}

		//		users[index] = u

	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.ConvertToLang(update.Message.From.LanguageCode, text.UnknownCommand),
		})
	}
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button(text.ConvertToLang(update.Message.From.LanguageCode, text.Join), []byte("1"), onJoinLobbySelect).
		Row().
		Button(text.ConvertToLang(update.Message.From.LanguageCode, text.Create), []byte("2"), onCreateLobbySelect)

	newUser := game.User{
		ChatID: update.Message.Chat.ID,
		Name: bot.EscapeMarkdown(update.Message.From.FirstName) + " " +
			bot.EscapeMarkdown(update.Message.From.LastName),
		Lang: update.Message.From.LanguageCode,
	}

	users.Append(newUser)

	log.Printf("New user: {Name: %s, ChatID: %d} is added", newUser.Name, newUser.ChatID)

	newUser.SendReplayMarkup(ctx, b, kb, text.StartCommandF, newUser.Name)
}

func GameStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	currentUser, exists := users.GetUser(update.Message.Chat.ID)
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.ConvertToLang(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		return
	}

	lobbies.Mut.RLock()
	lob, exists := lobbies.L[currentUser.LobbyKey]
	lobbies.Mut.RUnlock()

	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.ConvertToLang(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		return
	}

	if lob.GameType == text.GameNotSelected {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.ConvertToLang(update.Message.From.LanguageCode, text.CantStartGame),
		})
		return
	}

	if len(lob.Members) < game.MINIMUM_MEMBERS_FOR_MAFIA {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.ConvertToLang(update.Message.From.LanguageCode, text.AtLeastMembersF, game.MINIMUM_MEMBERS_FOR_MAFIA),
		})
		return
	}

	for _, member := range lob.Members {
		member.SendMessage(ctx, b, text.GameStarted)
	}

	lob.IsStarted = true

	lobbies.Mut.Lock()
	lobbies.L[currentUser.LobbyKey] = lob
	lobbies.Mut.Unlock()

	// go lobby.StartGame(ctx, b) ???

	lob.StartGame(ctx, b)
}
