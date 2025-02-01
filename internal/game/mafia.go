package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	//	"time"

	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

const (
	Civilian = iota
	Mafia
	Detective
	Doctor
)

type MafiaPlayer struct {
	lang		*string

	votes		int32

	role		int
	isAlive		bool
	isHealed	bool
}

func (m *MafiaPlayer) GetInfo() string {
	return text.GetConvertToLang(*m.lang, text.RoleF, m.getRole())
}

func (m MafiaPlayer) getRole() string {
	en := [...]string{
		"Civilian",
		"Mafia",
		"Detective",
		"Doctor",
	}

	ru := [...]string{
		"Мирный житель",
		"Мафиа",
		"Детектив",
		"Доктор",
	}

	switch *m.lang {
	case "en":
		return en[m.role]
	case "ru":
		return ru[m.role]
	default:
		return en[m.role]
	}
}

func getRole(lang string, role int) string {
	return MafiaPlayer{lang: &lang, role: role}.getRole()
}

type MafiaGame struct {
	isStarted	*bool

	members		*Users

	victim		*User

	mafias		UsersRef
	detectives	UsersRef
	doctors		UsersRef
}

func (m *MafiaGame) sendAll(ctx context.Context, b Bot, key int, mafiaKey int) {
	for _, member := range *m.members {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: member.ChatID,
			Text: text.GetConvertToLang(member.Lang, key, getRole(member.Lang, mafiaKey)),
		})
	}
}

func (m *MafiaGame) startGame(ctx context.Context, b Bot) {
	var wg sync.WaitGroup

	countOfMembers := len(*m.members)

	m.fillRoles(ctx, b)

	for *m.isStarted {
		m.members.sendAll(ctx, b, text.NightFalls)

		time.Sleep(2 * time.Second)

		// Mafia step

		wg.Add(1)

		m.sendAll(ctx, b, text.IsWakingUpF, Mafia)

		m.mafiaStep(ctx, b, &wg)

		// for i := range m.mafias {
		// 	m.mafias[i].SendReplayMarkup(ctx, b,
		// 		m.getMafiaKeyboard(b, &wg),
		// 		text.MakeChoiceF, Mafia,
		// 		m.mafias[i].Name,
		// 	)
		// }

		wg.Wait()

		m.sendAll(ctx, b, text.FallAsleepF, Mafia)

		time.Sleep(2 * time.Second)

		// Doctor step

		m.sendAll(ctx, b, text.IsWakingUpF, Doctor)

		time.Sleep(2 * time.Second)

		if !m.doctorsIsDead() {
			wg.Add(1)

			for i := range m.doctors {
				m.doctors[i].SendReplayMarkup(ctx, b,
					m.getDoctorKeyboard(b, &wg),
					text.MakeChoiceF, Doctor,
					m.doctors[i].Name,
				)
			}

			wg.Wait()
		}

		m.members.sendAll(ctx, b, text.FallAsleepF, Doctor)

		time.Sleep(2 * time.Second)

		// Detective step

		m.members.sendAll(ctx, b, text.IsWakingUpF, Detective)

		time.Sleep(2 * time.Second)

		if !m.detectivesIsDead() {
			wg.Add(1)

			for i := range m.detectives {
				m.detectives[i].SendReplayMarkup(ctx, b,
					m.getDetectiveKeyboard(b, &wg),
					text.MakeChoiceF, Detective,
					m.detectives[i].Name,
				)
			}
			
			wg.Wait()
		}

		m.sendAll(ctx, b, text.FallAsleepF, Detective)

		time.Sleep(2 * time.Second)

		m.members.sendAll(ctx, b, text.DayIsComing)

		if m.victim.player.(*MafiaPlayer).isHealed {
			m.members.sendAll(ctx, b, text.MafiaFailed)
			m.victim.player.(*MafiaPlayer).isHealed = false
			m.victim = nil
		} else {
			m.victim.player.(*MafiaPlayer).isAlive = false
			m.victim.player.(*MafiaPlayer).isHealed = false

			m.victim = nil
			
			m.members.sendAll(ctx, b, text.MafiaSuccessF, m.victim.Name)
		}

		// Voting session

		for {
			wg.Add(countOfMembers)

			m.runVote(ctx, b, &wg)

			wg.Wait()

			if !*m.isStarted {
				return
			}

			var voteResults string

			for _, mPlayer := range *m.members {
				if voteResults != "" {
					voteResults += "\n"
				}
				voteResults += mPlayer.Name + ": " + fmt.Sprintf("%d", mPlayer.player.(*MafiaPlayer).votes)
			}

			m.members.sendAll(ctx, b, text.VotingResultsF, voteResults)

			voteFirstMax, voteSecondMax := m.getTwoMaxVotes()

			if voteFirstMax.player.(*MafiaPlayer).votes == voteSecondMax.player.(*MafiaPlayer).votes {
				m.clearVotes()
				m.members.sendAll(ctx, b, text.VotesAreEqual)

			} else {
				m.members.sendAll(ctx, b, text.VoteKickF, voteFirstMax.Name, voteFirstMax.player.(*MafiaPlayer).role)
				m.kick(voteFirstMax)
				break
			}
		}

		if m.mafiaIsDead() {
			m.members.sendAll(ctx, b, text.CiviliansWon)
			break
		}

		if m.civiliansIsDead() {
			m.members.sendAll(ctx, b, text.MafiaWon)
			break
		}
	}
}

func (m *MafiaGame) fillRoles(ctx context.Context, b Bot) {
	size := len(*m.members)
	mSize := 1

	if size >= 6 {
		mSize = 2
	}

	for i := range mSize {
		m.mafias = append(m.mafias, &(*m.members)[rand.Intn(size)])
		m.mafias[i].player.(*MafiaPlayer).role = Mafia
	}

	for i := range 1 {
		m.detectives = append(m.detectives, &(*m.members)[rand.Intn(size)])
		m.detectives[i].player.(*MafiaPlayer).role = Detective
	}

	for i := range 1 {
		m.doctors = append(m.doctors, &(*m.members)[rand.Intn(size)])
		m.doctors[i].player.(*MafiaPlayer).role = Doctor
	}

	for i := range *m.members {
		(*m.members)[i].player.(*MafiaPlayer).isAlive = true
		(*m.members)[i].SendMessage(ctx, b, text.Default, (*m.members)[i].player.GetInfo())
	}
}

func (m *MafiaGame) getDoctorKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {
	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		player, ok := m.members.getMember(int64(ChID))
		if !ok {
			log.Panic("Can't find member in lobby")
		}

//		player := &(*m.members)[m.findMember(User{ChatID: int64(ChID)})]

		player.player.(*MafiaPlayer).isHealed = true

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.HealedF, player.Name),
		})

	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.members {
		if player.player.(*MafiaPlayer).role != Doctor {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onVictimSelect)
		}
	}

	return kb
}

func (m *MafiaGame) getDetectiveKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {
	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		victim, ok := m.members.getMember(int64(ChID))
		if !ok {
			log.Print("Can't find member in lobby")
			
		}
//		victim := (*m.members)[m.findMember(User{ChatID: int64(ChID)})]
		if victim.player.(*MafiaPlayer).role == Mafia {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.IsMafiaF, victim.Name),
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.IsNotMafiaF, victim.Name),
			})
		}
	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.members {
		if player.player.(*MafiaPlayer).role != Detective {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onVictimSelect)
		}
	}

	return kb
}

func (m *MafiaGame) getMafiaKeyboard(b Bot, wg *sync.WaitGroup) *inline.Keyboard {
	onVictimSugest := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte)  {
		onYes := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()
			

			ChID, _ := strconv.Atoi(string(data))
			m.victim = &(*m.members)[m.members.findMember(User{ChatID: int64(ChID)})]
			this := (*m.members)[m.members.findMember(User{ChatID: mes.Message.Chat.ID})]

			for _, maf := range m.mafias {
				if maf.ChatID != mes.Message.Chat.ID {
					maf.SendMessage(ctx, b, text.AcceptedF, this.Name)
					break
				} 
			}
		}

		onNo := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()
			
			this := (*m.members)[m.members.findMember(User{ChatID: mes.Message.Chat.ID})]

			for _, maf := range m.mafias {
				if maf.ChatID != mes.Message.Chat.ID {
					maf.SendMessage(ctx, b, text.DeclinedF, this.Name)
					break
				} 
			}
		}

		ChID, _ := strconv.Atoi(string(data))
		probVictim := (*m.members)[m.members.findMember(User{ChatID: int64(ChID)})]
		whoSugest := (*m.members)[m.members.findMember(User{ChatID: mes.Message.Chat.ID})]

		for _, maf := range m.mafias {
			if mes.Message.Chat.ID != maf.ChatID {
				kb := inline.New(b).
						Row().
							Button(text.GetConvertToLang(mes.Message.From.LanguageCode, text.Yes), []byte(fmt.Sprintf("%d", probVictim.ChatID)), onYes).
							Button(text.GetConvertToLang(mes.Message.From.LanguageCode, text.No), []byte(fmt.Sprintf("%d", probVictim.ChatID)), onNo)

				maf.SendReplayMarkup(ctx, b, kb, text.SugestToKillF, whoSugest.Name, probVictim.Name)
			}
		}
	}

	onVictimSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		defer wg.Done()

		ChID, _ := strconv.Atoi(string(data))
		m.victim = &(*m.members)[m.members.findMember(User{ChatID: int64(ChID)})]

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.ChosenToKillF, m.victim.Name),
		})
	}

	var onSelect inline.OnSelect

	if len(m.mafias) > 1 {
		onSelect = onVictimSugest
	} else {
		onSelect = onVictimSelect
	}

	kb := inline.New(b.(*bot.Bot))

	for _, player := range *m.members {
		if player.player.(*MafiaPlayer).role != Mafia {
			kb.Row().Button(player.Name, []byte(fmt.Sprintf("%d", player.ChatID)), onSelect)
		}
	}

	return kb
}

func (m *MafiaGame) mafiaStep(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	defer wg.Done()

	var innerWg sync.WaitGroup

	for m.victim == nil || (m.victim != nil && !m.victim.player.(*MafiaPlayer).isAlive) {
		innerWg.Add(1)

		for i := range m.mafias {
			m.mafias[i].SendReplayMarkup(ctx, b,
				m.getMafiaKeyboard(b, &innerWg),
				text.MakeChoiceF, m.mafias[i].player.(*MafiaPlayer).getRole(),
				m.mafias[i].Name,
			)
		}

		innerWg.Wait()
	}
}

func (m *MafiaGame) runVote(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for i := range *m.members {
		onVoteSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			ChID, _ := strconv.Atoi(string(data))
			player, ok := m.members.getMember(int64(ChID))
			if !ok {
				log.Print("Can't find member in lobby")
				*m.isStarted = false
				m.members.sendAll(ctx, b, text.GameStoped)
				return
			}
			atomic.AddInt32(&player.player.(*MafiaPlayer).votes, 1)
			player.player.(*MafiaPlayer).votes++

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text.GetConvertToLang(mes.Message.From.LanguageCode, text.VotedF, player.Name),
			})
		}

		kb := inline.New(b.(*bot.Bot))

		for j := range *m.members {
			if (*m.members)[i].ChatID == (*m.members)[j].ChatID {
				continue
			}
			kb.Row().Button((*m.members)[j].Name, []byte(fmt.Sprintf("%d", (*m.members)[j].ChatID)), onVoteSelect)
		}
		
		(*m.members)[i].SendReplayMarkup(ctx, b, kb, text.Voting)
	}
}

func (m *MafiaGame) mafiaIsDead() bool {
	for i := range m.mafias {
		if m.mafias[i].player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) civiliansIsDead() bool {
	for i := range *m.members {
		if (*m.members)[i].player.(*MafiaPlayer).role != Mafia && (*m.members)[i].player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) detectivesIsDead() bool {
	for i := range m.detectives {
		if m.detectives[i].player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) doctorsIsDead() bool {
	for i := range m.doctors {
		if m.doctors[i].player.(*MafiaPlayer).isAlive {
			return false
		}
	}
	return true
}

func (m *MafiaGame) getTwoMaxVotes() (User, User) {
	maxFirst := User{player: &MafiaPlayer{votes: 0}}
	maxSecond := maxFirst

	for _, player := range *m.members {
		if maxFirst.player.(*MafiaPlayer).votes < player.player.(*MafiaPlayer).votes {
			maxSecond = maxFirst
			maxFirst = player
		} else if maxSecond.player.(*MafiaPlayer).votes < player.player.(*MafiaPlayer).votes {
			maxSecond = player
		}
	}

	return maxFirst, maxSecond
}

func (m *MafiaGame) clearVotes() {
	for i := range *m.members {
		(*m.members)[i].player.(*MafiaPlayer).votes = 0
	}
}

func (m *MafiaGame) kick(u User) int {
	index := m.members.findMember(u)
	(*m.members)[index].player.(*MafiaPlayer).isAlive = false

	return (*m.members)[index].player.(*MafiaPlayer).role
}