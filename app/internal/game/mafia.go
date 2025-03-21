package game

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	//	"time"

	"gffbot/internal/logger"
	"gffbot/internal/storage"
	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"go.uber.org/zap"
)

type MafiaPlayer struct {
	Lang *string

	votes int32

	role     int
	isAlive  bool
	isHealed bool
}

func (m *MafiaPlayer) Info() string {
	return text.Convert(*m.Lang, text.RoleF, m.getRole())
}

func (m MafiaPlayer) getRole() string {
	return text.Convert(*m.Lang, m.role)
}

type MafiaGame struct {
	IsStarted *bool

	mut     sync.Mutex
	Members *Users

	victim *User

	mafias     UsersRef
	detectives UsersRef
	doctors    UsersRef
}

func (m *MafiaGame) StartGame(ctx context.Context, b Bot, repo *storage.Repository) {
	var wg sync.WaitGroup

	countOfAliveMembers := len(*m.Members)

	m.fillRoles(ctx, b)

	for *m.IsStarted {
		m.Members.SendAll(ctx, b, text.NightFalls)

		time.Sleep(2 * time.Second)

		// Mafia step

		wg.Add(1) // check !!!

		m.Members.SendAll(ctx, b, text.IsWakingUpF, text.Mafia)

		m.mafiaStep(ctx, b, &wg)

		wg.Wait()

		m.Members.SendAll(ctx, b, text.FallAsleepF, text.Mafia)

		time.Sleep(2 * time.Second)

		// Doctor step

		m.Members.SendAll(ctx, b, text.IsWakingUpF, text.Doctor)

		time.Sleep(2 * time.Second)

		if !m.doctorsIsDead() {
			wg.Add(1)

			for i := range m.doctors {
				m.doctors[i].SendReplayMarkup(ctx, b,
					m.doctorKeyboard(b, &wg),
					text.MakeChoiceF, text.Doctor,
					m.doctors[i].Name,
				)
			}

			wg.Wait()
		}

		m.Members.SendAll(ctx, b, text.FallAsleepF, text.Doctor)

		time.Sleep(2 * time.Second)

		// Detective step

		m.Members.SendAll(ctx, b, text.IsWakingUpF, text.Detective)

		time.Sleep(2 * time.Second)

		if !m.detectivesIsDead() {
			wg.Add(1)

			for i := range m.detectives {
				m.detectives[i].SendReplayMarkup(ctx, b,
					m.detectiveKeyboard(b, &wg),
					text.MakeChoiceF, text.Detective,
					m.detectives[i].Name,
				)
			}

			wg.Wait()
		}

		m.Members.SendAll(ctx, b, text.FallAsleepF, text.Detective)

		time.Sleep(2 * time.Second)

		m.Members.SendAll(ctx, b, text.DayIsComing)

		if m.victim.Player.(*MafiaPlayer).isHealed {
			m.Members.SendAll(ctx, b, text.MafiaFailed)
			m.victim.Player.(*MafiaPlayer).isHealed = false
			m.victim = nil
		} else {
			m.victim.Player.(*MafiaPlayer).isAlive = false
			m.victim.Player.(*MafiaPlayer).isHealed = false

			m.victim = nil

			m.Members.SendAll(ctx, b, text.MafiaSuccessF, m.victim.Name)

			countOfAliveMembers--
		}

		// Voting session

		for {
			wg.Add(countOfAliveMembers)

			m.runVote(ctx, b, &wg)

			wg.Wait()

			if !*m.IsStarted {
				return
			}

			var voteResults string

			for _, mPlayer := range *m.Members {
				if voteResults != "" {
					voteResults += "\n"
				}
				voteResults += mPlayer.Name + ": " + fmt.Sprintf("%d", mPlayer.Player.(*MafiaPlayer).votes)
			}

			m.Members.SendAll(ctx, b, text.VotingResultsF, voteResults)

			voteFirstMax, voteSecondMax := m.twoMaxVotes()

			if voteFirstMax.Player.(*MafiaPlayer).votes == voteSecondMax.Player.(*MafiaPlayer).votes {
				m.clearVotes()
				m.Members.SendAll(ctx, b, text.VotesAreEqual)

			} else {
				m.Members.SendAll(ctx, b, text.VoteKickF, voteFirstMax.Name, voteFirstMax.Player.(*MafiaPlayer).role)
				m.kick(voteFirstMax)
				countOfAliveMembers--
				break
			}
		}

		if m.unwinnable() {
			m.Members.SendAll(ctx, b, text.MafiaWon)
			m.loadResults(repo, true)
			break
		}

		if m.mafiaIsDead() {
			m.Members.SendAll(ctx, b, text.CiviliansWon)
			m.loadResults(repo, false)
			break
		}

		if m.civiliansIsDead() {
			m.Members.SendAll(ctx, b, text.MafiaWon)
			m.loadResults(repo, true)
			break
		}
	}
}

func (m *MafiaGame) fillRoles(ctx context.Context, b Bot) {
	size := len(*m.Members)
	mSize := 1
	tofill := 0

	if size >= 6 {
		mSize = 2
	}

	for i := range mSize {
		tofill = rand.IntN(size)

		for (*m.Members)[tofill].Player.(*MafiaPlayer).role != 0 {
			tofill = rand.IntN(size)
		}

		m.mafias = append(m.mafias, &(*m.Members)[tofill])
		m.mafias[i].Player.(*MafiaPlayer).role = text.Mafia
	}

	for i := range 1 {
		tofill = rand.IntN(size)

		for (*m.Members)[tofill].Player.(*MafiaPlayer).role != 0 {
			tofill = rand.IntN(size)
		}

		m.detectives = append(m.detectives, &(*m.Members)[tofill])
		m.detectives[i].Player.(*MafiaPlayer).role = text.Detective
	}

	for i := range 1 {
		tofill = rand.IntN(size)

		for (*m.Members)[tofill].Player.(*MafiaPlayer).role != 0 {
			tofill = rand.IntN(size)
		}

		m.doctors = append(m.doctors, &(*m.Members)[tofill])
		m.doctors[i].Player.(*MafiaPlayer).role = text.Doctor
	}

	for i := range *m.Members {
		(*m.Members)[i].Player.(*MafiaPlayer).isAlive = true
		(*m.Members)[i].SendMessage(ctx, b, text.Default, (*m.Members)[i].Player.Info())
	}
}

func (m *MafiaGame) doctorKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {

	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		player, ok := m.Members.GetMember(int64(ChID))
		if !ok {
			logger.Log.Info("can't find member in lobby", zap.Int64("chat_id", mes.Message.Chat.ID), zap.Int64("mising_chat_id", int64(ChID)))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text: text.Convert(mes.Message.From.LanguageCode, text.SomethingWentWrong),
			})
			return
		}

		player.Player.(*MafiaPlayer).isHealed = true

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.Convert(mes.Message.From.LanguageCode, text.HealedF, player.Name),
		})

	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.Members {
		if player.Player.(*MafiaPlayer).role != text.Doctor {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onVictimSelect)
		}
	}

	return kb
}

func (m *MafiaGame) detectiveKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {
	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		victim, ok := m.Members.GetMember(int64(ChID))
		if !ok {
			logger.Log.Info("can't find member in lobby", zap.Int64("chat_id", mes.Message.Chat.ID), zap.Int64("mising_chat_id", int64(ChID)))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text: text.Convert(mes.Message.From.LanguageCode, text.SomethingWentWrong),
			})
			return
		}

		if victim.Player.(*MafiaPlayer).role == text.Mafia {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.Convert(mes.Message.From.LanguageCode, text.IsMafiaF, victim.Name),
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.Convert(mes.Message.From.LanguageCode, text.IsNotMafiaF, victim.Name),
			})
		}
	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.Members {
		if player.Player.(*MafiaPlayer).role != text.Detective {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onVictimSelect)
		}
	}

	return kb
}

func (m *MafiaGame) mafiaKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {
	onVictimSugest := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		onYes := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			ChID, _ := strconv.Atoi(string(data))
			m.victim = &(*m.Members)[m.Members.FindMember(User{ChatID: int64(ChID)})]
			this := (*m.Members)[m.Members.FindMember(User{ChatID: mes.Message.Chat.ID})]

			for _, maf := range m.mafias {
				if maf.ChatID != mes.Message.Chat.ID {
					maf.SendMessage(ctx, b, text.AcceptedF, this.Name)
					break
				}
			}
		}

		onNo := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			this := (*m.Members)[m.Members.FindMember(User{ChatID: mes.Message.Chat.ID})]

			for _, maf := range m.mafias {
				if maf.ChatID != mes.Message.Chat.ID {
					maf.SendMessage(ctx, b, text.DeclinedF, this.Name)
					break
				}
			}
		}

		ChID, _ := strconv.Atoi(string(data))
		probVictim := (*m.Members)[m.Members.FindMember(User{ChatID: int64(ChID)})]
		whoSugest := (*m.Members)[m.Members.FindMember(User{ChatID: mes.Message.Chat.ID})]

		for _, maf := range m.mafias {
			if mes.Message.Chat.ID != maf.ChatID {
				kb := inline.New(b).
					Row().
					Button(text.Convert(mes.Message.From.LanguageCode, text.Yes), []byte(fmt.Sprintf("%d", probVictim.ChatID)), onYes).
					Button(text.Convert(mes.Message.From.LanguageCode, text.No), []byte(fmt.Sprintf("%d", probVictim.ChatID)), onNo)

				maf.SendReplayMarkup(ctx, b, kb, text.SugestToKillF, whoSugest.Name, probVictim.Name)
			}
		}
	}

	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		m.victim = &(*m.Members)[m.Members.FindMember(User{ChatID: int64(ChID)})]

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.Convert(mes.Message.From.LanguageCode, text.ChosenToKillF, m.victim.Name),
		})
	}

	var onSelect inline.OnSelect

	if len(m.mafias) > 1 {
		onSelect = onVictimSugest
	} else {
		onSelect = onVictimSelect
	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.Members {
		if player.Player.(*MafiaPlayer).role != text.Mafia {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onSelect)
		}
	}

	return kb
}

func (m *MafiaGame) mafiaStep(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	defer wg.Done()

	var innerWg sync.WaitGroup

	for m.victim == nil || (m.victim != nil && !m.victim.Player.(*MafiaPlayer).isAlive) {
		innerWg.Add(1)

		for i := range m.mafias {
			m.mafias[i].SendReplayMarkup(ctx, b,
				m.mafiaKeyboard(b, &innerWg),
				text.MakeChoiceF, m.mafias[i].Player.(*MafiaPlayer).getRole(),
				m.mafias[i].Name,
			)
		}

		innerWg.Wait()
	}
}

func (m *MafiaGame) runVote(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for i := range *m.Members {
		onVoteSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			ChID, _ := strconv.Atoi(fmt.Sprintf("%d", data))
			player, ok := m.Members.GetMember(int64(ChID))
			if !ok {
				logger.Log.Info("can't find member in lobby", zap.Int64("chat_id", mes.Message.Chat.ID), zap.Int64("mising_chat_id", int64(ChID)))
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text: text.Convert(mes.Message.From.LanguageCode, text.SomethingWentWrong),
				})
				*m.IsStarted = false
				m.Members.SendAll(ctx, b, text.GameStoped)
				return
			}
			atomic.AddInt32(&player.Player.(*MafiaPlayer).votes, 1)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.Convert(mes.Message.From.LanguageCode, text.VotedF, player.Name),
			})
		}

		kb := inline.New(b.(*bot.Bot))

		for j := range *m.Members {
			if (*m.Members)[i].ChatID == (*m.Members)[j].ChatID {
				continue
			}
			kb.Row().Button((*m.Members)[j].Name, []byte(fmt.Sprintf("%d", (*m.Members)[j].ChatID)), onVoteSelect)
		}

		(*m.Members)[i].SendReplayMarkup(ctx, b, kb, text.Voting)
	}
}

func (m *MafiaGame) mafiaIsDead() bool {
	for i := range m.mafias {
		if m.mafias[i].Player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) civiliansIsDead() bool {
	for i := range *m.Members {
		if (*m.Members)[i].Player.(*MafiaPlayer).role != text.Mafia && (*m.Members)[i].Player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) detectivesIsDead() bool {
	for i := range m.detectives {
		if m.detectives[i].Player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) unwinnable() bool {
	mafiaCount := 0
    civilianCount := 0

    for _, player := range *m.Members {
        if player.Player.(*MafiaPlayer).role == text.Mafia {
            mafiaCount++
        } else {
            civilianCount++
        }
    }

    return mafiaCount >= civilianCount
}

func (m *MafiaGame) doctorsIsDead() bool {
	for i := range m.doctors {
		if m.doctors[i].Player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) twoMaxVotes() (User, User) {
	maxFirst := User{Player: &MafiaPlayer{votes: 0}}
	maxSecond := maxFirst

	for _, player := range *m.Members {
		if maxFirst.Player.(*MafiaPlayer).votes < player.Player.(*MafiaPlayer).votes {
			maxSecond = maxFirst
			maxFirst = player
		} else if maxSecond.Player.(*MafiaPlayer).votes < player.Player.(*MafiaPlayer).votes {
			maxSecond = player
		}
	}

	return maxFirst, maxSecond
}

func (m *MafiaGame) clearVotes() {
	for i := range *m.Members {
		(*m.Members)[i].Player.(*MafiaPlayer).votes = 0
	}
}

func (m *MafiaGame) kick(u User) int {
	index := m.Members.FindMember(u)
	(*m.Members)[index].Player.(*MafiaPlayer).isAlive = false

	return (*m.Members)[index].Player.(*MafiaPlayer).role
}

func (m *MafiaGame) loadResults(repo *storage.Repository, mafiaWon bool) {
	var wg sync.WaitGroup
	wg.Add(len(*m.Members))

	for _, member := range *m.Members {
		go func(user User) {
			defer wg.Done()
			if user.Player.(*MafiaPlayer).role == text.Mafia {
				err := repo.UpdateUserStatistic(user.ID, mafiaWon)
				if err != nil {
					logger.Log.Error("can't update user statistic", zap.Int64("user_id", user.ID), zap.Error(err))
				}
			} else {
				err := repo.UpdateUserStatistic(user.ID, !mafiaWon)
				if err != nil {
					logger.Log.Error("can't update user statistic", zap.Int64("user_id", user.ID), zap.Error(err))
				}
			}
		}(member)
	}

	wg.Wait()
}