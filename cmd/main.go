package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"gffbot/internal/game"
	"gffbot/internal/text"

	"github.com/joho/godotenv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

var users []game.User
var lobbies map[string]game.Lobby

func findUserInData(u game.User) (int, bool) {
	for i, user := range users {
		if u.ChatID == user.ChatID {
			return i, true
		}
	}
	return -1, false
}

func createLobbyKey() string {
	key := make([]byte, 4)

	key[0] = text.LettersBytes[rand.Intn(len(text.LettersBytes))]
	key[2] = text.LettersBytes[rand.Intn(len(text.LettersBytes))]
	key[1] = text.DigitsBytes[rand.Intn(len(text.DigitsBytes))]
	key[3] = text.DigitsBytes[rand.Intn(len(text.DigitsBytes))]

	return string(key)
}

func lastLobbyID() int64 {
	var m int64 = -1

	for _, lobby := range lobbies {
		m = max(m, lobby.ID)
	}

	return m
}

func init() {
	lobbies = make(map[string]game.Lobby)

	err := godotenv.Load("../env/.env")
	if err != nil {
		log.Panic(err)
	}

	rand.Seed(time.Now().UnixNano())
}

func main() {
	TOKEN := "BOT_API_TOKEN"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	gffbot, err := bot.New(os.Getenv(TOKEN), opts...)
	if err != nil {
		log.Panic(err)
	}

	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	gffbot.RegisterHandler(bot.HandlerTypeMessageText, "/game_start", bot.MatchTypeExact, gameStartHandler)

	gffbot.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	index, exists := findUserInData(game.User{ChatID: update.Message.Chat.ID})
	if exists && users[index].SendingKey {
		u := users[index]

		if u.LobbyID != 0 || u.LobbyKey != "" {
			u.SendMessage(ctx, b, text.AlreadyInLobbyF, u.LobbyKey)
			return
		}

		key := update.Message.Text

		if lobby, exists := lobbies[key]; exists {

			if lobby.IsStarted {
				u.SendMessage(ctx, b, text.LobbyGameIsStarted)
				return
			}

			// Поменять

			memebersList := []string{u.Name}

			for _, memeber := range lobby.Members {
				memebersList = append(memebersList, memeber.Name)
			}

			newList := strings.Join(memebersList, "\n")

			// ==//==

			for _, member := range lobby.Members {
				member.SendMessage(ctx, b, text.PlayerJoinedLobbyF, u.Name, newList)
			}

			u.SendingKey = false
			u.LobbyKey = key
			u.LobbyID = lobby.ID
			users[index] = u

			lobby.Members = append(lobby.Members, u)

			u.SendMessage(ctx, b, text.PlayerJoinedLobbyF, text.PlayerJoinedLobbyF, u.GetText(text.You), newList)
		} else {
			u.SendMessage(ctx, b, text.LobbyNotExists)
			return
		}

		//		users[index] = u

	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.GetConvertToLang(update.Message.From.LanguageCode, text.UnknownCommand),
		})
	}
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button(text.GetConvertToLang(update.Message.From.LanguageCode, text.Join), []byte("1"), onJoinLobbySelect).
		Row().
		Button(text.GetConvertToLang(update.Message.From.LanguageCode, text.Create), []byte("2"), onCreateLobbySelect)

	newUser := game.User{
		ChatID:	update.Message.Chat.ID,
		Name:	bot.EscapeMarkdown(update.Message.From.FirstName) + " " +
					bot.EscapeMarkdown(update.Message.From.LastName),
		Lang:	update.Message.From.LanguageCode,
	}

	users = append(users, newUser)

	log.Printf("New user: {Name: %s, ChatID: %d} is added", newUser.Name, newUser.ChatID)

	newUser.SendReplayMarkup(ctx, b, kb, text.StartCommandF, newUser.Name)
}

func gameStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	index, exists := findUserInData(game.User{ChatID: update.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.GetConvertToLang(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		return
	}

	lobby, exists := lobbies[users[index].LobbyKey]
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.GetConvertToLang(update.Message.From.LanguageCode, text.SomethingWentWrong),
		})
		return
	}

	if lobby.GameType == text.GameNotSelected {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.GetConvertToLang(update.Message.From.LanguageCode, text.CantStartGame),
		})
		return
	}

	if len(lobby.Members) < game.MINIMUM_MEMBERS_FOR_MAFIA {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.GetConvertToLang(update.Message.From.LanguageCode, text.AtLeastMembersF, game.MINIMUM_MEMBERS_FOR_MAFIA),
		})
		return
	}

	for _, member := range lobby.Members {
		member.SendMessage(ctx, b, text.GameStarted)
	}

	lobby.IsStarted = true

	lobbies[users[index].LobbyKey] = lobby

	// go lobby.StartGame(ctx, b) ???

	lobby.StartGame(ctx, b)
}

func onJoinLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	users[index].SendingKey = true

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.SendKey),
	})
}

func onGameSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	gameType, _ := strconv.Atoi(string(data))

	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	user := users[index]

	lobby, exists := lobbies[user.LobbyKey]
	if !exists {
		user.SendMessage(ctx, b, text.SomethingWentWrong)
		log.Printf("key: %s", user.LobbyKey)
		return
	}

	lobby.GameType = gameType
	lobbies[user.LobbyKey] = lobby

	user.SendMessage(ctx, b, text.GameChosenF, gameType)
}

func onCreateLobbySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	index, exists := findUserInData(game.User{ChatID: mes.Message.Chat.ID})
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CantFindUser),
		})
		return
	}

	user := users[index]

	var key string

	if len(lobbies) != 0 {
		for i := range 10 {
			key = createLobbyKey()

			if _, exists := lobbies[key]; exists {
				break
			} else if i == 9 {
				log.Printf("Somefting went wrong on creating new lobby. Current count of lobbies is: %d", len(lobbies))
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.CreatingLobbyError),
				})
			}
		}
	} else {
		key = createLobbyKey()
	}

	newLobby := game.Lobby{
		LeaderID:  user.ChatID,
		GameType:  text.GameNotSelected,
		IsStarted: false,
		Members:   []game.User{user},
	}

	lobbies[key] = newLobby

	user.LobbyKey = key
	users[index] = user

	kb := inline.New(b).
		Row().
		Button(game.GetGame(mes.Message.From.LanguageCode, game.GMafia), []byte(fmt.Sprintf("%d", game.GMafia)), onGameSelect).
		Row().
		Button(game.GetGame(mes.Message.From.LanguageCode, game.GBunker), []byte(fmt.Sprintf("%d", game.GBunker)), onGameSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        text.GetConvertToLang(mes.Message.From.LanguageCode, text.KeyCreatedF, key),
		ReplyMarkup: kb,
	})

	log.Printf("New lobby {ID: %d, key: %s} added", newLobby.ID, key)
}
