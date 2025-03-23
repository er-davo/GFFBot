package game

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"sync/atomic"

	"gffbot/internal/logger"
	"gffbot/internal/storage"
	"gffbot/internal/text"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/paginator"
	"go.uber.org/zap"
)

type BunkerPlayerFeature struct {
	info     string
	isHidden bool
}

func (bf *BunkerPlayerFeature) toString() string {
	if bf.isHidden {
		return "[скрыто]  " + bf.info
	} else {
		return "[открыто] " + bf.info
	}
}

func (bf *BunkerPlayerFeature) view() string {
	if bf.isHidden {
		return "скрыто"
	}
	return bf.info
}

type BunkerPlayer struct {
	Lang *string

	profession string

	biologicalParams BunkerPlayerFeature
	healthStatus     BunkerPlayerFeature
	hobby            BunkerPlayerFeature
	phobia           BunkerPlayerFeature
	character        BunkerPlayerFeature
	additionalInfo   BunkerPlayerFeature
	skill            BunkerPlayerFeature
	knowledge        BunkerPlayerFeature
	baggage          BunkerPlayerFeature

	actionCard    any
	conditionCard any

	votes    int32
	isKicked bool
}

// Still TODO
func (bp *BunkerPlayer) fill() {

	//	luckyFeatures := 4

	// Biological Params

	age := 12 + rand.IntN(79)
	bp.biologicalParams.info = fmt.Sprintf("%d", age)

	if age%10 == 1 {
		bp.biologicalParams.info += " год"
	} else if age%10 < 5 {
		bp.biologicalParams.info += " года"
	} else {
		bp.biologicalParams.info += " лет"
	}

	sex := [2]string{"Мужчина", "Женщина"}[rand.IntN(2)]
	bp.biologicalParams.info += "/" + sex

	if sex == "Женщина" && rand.IntN(25) == 0 {
		bp.biologicalParams.info += "/" + "беременна"
	} else {
		if rand.IntN(5) != 0 {
			bp.biologicalParams.info += "/" + text.BiologicalParamsRu[0]
		} else {
			bp.biologicalParams.info += "/" + text.BiologicalParamsRu[rand.IntN(text.BiologicalParamsLen)]
		}
	}

	// Profession

	switch sex {
	case "Мужчина":
		bp.profession = text.ProfessionsRu[rand.IntN(text.ProfessionsMale)]
	case "Женщина":
		bp.profession = text.ProfessionsRu[rand.IntN(text.ProfessionsLen-text.ProfessionsFemale)+text.ProfessionsFemale]
	}

	// Health status

	switch sex {
	case "Мужчина":
		bp.healthStatus.info = text.HealthStatusRu[rand.IntN(text.HealthStatusMale)]
	case "Женщина":
		bp.healthStatus.info = text.HealthStatusRu[rand.IntN(text.HealthStatusLen-text.HealthStatusFemale)+text.HealthStatusFemale]
	}

	// Hobbies

	bp.hobby.info = text.HobbiesRu[rand.IntN(text.HobbiesLen)]

	// Phobia

	bp.phobia.info = text.PhobiasRu[rand.IntN(text.PhopiaLen)]

	// Character

	bp.character.info = text.CharacterRu[rand.IntN(text.CharacterLen)]

	// Skill

	bp.skill.info = text.SkillsRu[rand.IntN(text.SkillsLen)]

	// Baggage

	bp.baggage.info = text.BaggageRu[rand.IntN(text.BaggageLen)]

	bp.biologicalParams.isHidden = true
	bp.healthStatus.isHidden = true
	bp.hobby.isHidden = true
	bp.phobia.isHidden = true
	bp.character.isHidden = true
	bp.additionalInfo.isHidden = true
	bp.skill.isHidden = true
	bp.knowledge.isHidden = true
	bp.baggage.isHidden = true
}

func (bp *BunkerPlayer) Info() string {
	var info string

	if bp.isKicked {
		info = "(Выгнан)\n"
	}

	info += text.Profession + ": " + bp.profession + "\n"

	return info + text.BoilogicalParams + ": " +
		bp.biologicalParams.toString() +
		"\n" + text.HealthStatus + ": " +
		bp.healthStatus.toString() +
		"\n" + text.Hobby + ": " +
		bp.hobby.toString() +
		"\n" + text.Phopia + ": " +
		bp.phobia.toString() +
		"\n" + text.Character + ": " +
		bp.character.toString() +
		"\n" + text.Skill + ": " +
		bp.skill.toString() +
		"\n" + text.Knowledge + ": " +
		bp.knowledge.toString() +
		"\n" + text.Baggage + ": " +
		bp.baggage.toString() +
		"\n\n" + text.ActionCard + ": " +
		"TODO" +
		"\n" + text.ConditionCard + ": " +
		"TODO"
}

func (bp *BunkerPlayer) View() string {
	var info string

	if bp.isKicked {
		info = "(Выгнан)\n"
	}

	info += text.Profession + ": " + bp.profession + "\n"

	return info + text.BoilogicalParams + ": " +
		bp.biologicalParams.view() +
		"\n" + text.HealthStatus + ": " +
		bp.healthStatus.view() +
		"\n" + text.Hobby + ": " +
		bp.hobby.view() +
		"\n" + text.Phopia + ": " +
		bp.phobia.view() +
		"\n" + text.Character + ": " +
		bp.character.view() +
		"\n" + text.Skill + ": " +
		bp.skill.view() +
		"\n" + text.Knowledge + ": " +
		bp.knowledge.view() +
		"\n" + text.Baggage + ": " +
		bp.baggage.view()
}

func (bp *BunkerPlayer) Keyboard(b Bot, onSelect inline.OnSelect) *inline.Keyboard {
	kb := inline.New(b.(*bot.Bot))

	if bp.biologicalParams.isHidden {
		kb.Row().Button(text.BoilogicalParams, []byte("1"), onSelect)
	}
	if bp.healthStatus.isHidden {
		kb.Row().Button(text.HealthStatus, []byte("2"), onSelect)
	}
	if bp.hobby.isHidden {
		kb.Row().Button(text.Hobby, []byte("3"), onSelect)
	}
	if bp.phobia.isHidden {
		kb.Row().Button(text.Phopia, []byte("4"), onSelect)
	}
	if bp.character.isHidden {
		kb.Row().Button(text.Character, []byte("5"), onSelect)
	}
	if bp.skill.isHidden {
		kb.Row().Button(text.Skill, []byte("6"), onSelect)
	}
	if bp.knowledge.isHidden {
		kb.Row().Button(text.Knowledge, []byte("7"), onSelect)
	}
	if bp.baggage.isHidden {
		kb.Row().Button(text.Baggage, []byte("8"), onSelect)
	}

	return kb
}

type BunkerGame struct {
	disastre  string
	IsStarted *bool
	Members   *Users
}

func (bg *BunkerGame) send(ctx context.Context, b Bot, msg string, a ...any) {
	var wg sync.WaitGroup
	wg.Add(len(*bg.Members))

	for _, player := range *bg.Members {
		go func(user User) {
			defer wg.Done()
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: user.ChatID,
				Text:   fmt.Sprintf(msg, a...),
			})
		}(player)
	}

	wg.Wait()
}

var toKickData = [...][5]int{
	{0, 0, 0, 1, 2},  // 3
	{0, 0, 1, 1, 2},  // 4 players
	{0, 1, 1, 1, 2},  // 5
	{0, 1, 1, 1, 3},  // 6
	{1, 1, 1, 1, 3},  // 7
	{1, 1, 1, 1, 4},  // 8
	{1, 1, 1, 2, 4},  // 9
	{1, 1, 1, 2, 5},  // 10
	{1, 1, 2, 2, 5},  // 11
	{1, 1, 2, 2, 6},  // 12
	{1, 2, 2, 2, 6},  // 13
	{1, 2, 2, 2, 7},  // 14
	{2, 2, 2, 2, 7},  // 15
	{2, 2, 2, 2, 8},  // 16
	{2, 2, 2, 3, 8},  // 17
	{2, 2, 2, 3, 9},  // 18
	{2, 2, 3, 3, 9},  // 19
	{2, 2, 3, 3, 10}, // 20
	{2, 3, 3, 3, 10}, // 21
	{2, 3, 3, 3, 11}, // 22
}

func (bg *BunkerGame) StartGame(ctx context.Context, b Bot, repo *storage.Repository) {
	// fill roles and disastre

	var wg sync.WaitGroup

	countOfAliveMembers := len(*bg.Members)
	tokick := toKickData[countOfAliveMembers-3]

	bg.fillFeatures()
	bg.send(ctx, b, bg.disastre+"\nКоличесвто мест: %d", tokick[4])

	for step := 0; step < 4; step++ {
		wg.Add(countOfAliveMembers)

		bg.openHiddenFeatures(ctx, b, &wg)

		wg.Wait()
		wg.Add(countOfAliveMembers)

		bg.openHiddenFeatures(ctx, b, &wg)

		wg.Wait()
		wg.Add(countOfAliveMembers)

		bg.sendAllInfo(ctx, b, &wg)

		wg.Wait()

		for range tokick[step] {
			for {
				wg.Add(countOfAliveMembers)

				bg.runVote(ctx, b, &wg)

				wg.Wait()

				voteMessage := "Результаты голосования:\n" + bg.voteResults()
				bg.send(ctx, b, voteMessage)

				voteFirstMax, voteSecondMax := bg.twoMaxVotes()

				if voteFirstMax.Player.(*BunkerPlayer).votes == voteSecondMax.Player.(*BunkerPlayer).votes {
					bg.clearVotes()
					bg.send(ctx, b, text.VotesAreEqualRu)
				} else {
					bg.send(ctx, b, text.KickedF, voteFirstMax.Name)
					// change
					bg.Members.SendAll(ctx, b, text.VoteKickF, voteFirstMax.Name, voteFirstMax.Player.(*BunkerPlayer))
					bg.kick(voteFirstMax)
					countOfAliveMembers--
					break
				}
			}
		}
	}

	// ai summary

	bg.send(ctx, b, text.GameEnd)
	bg.loadResults(repo)
}

func (bg *BunkerGame) fillFeatures() {
	bg.disastre = text.ApocalypseExamples[rand.IntN(text.ApocalypseExamplesLen)]

	var wg sync.WaitGroup
	wg.Add(len(*bg.Members))

	for iter := range *bg.Members {
		go func(i int) {
			defer wg.Done()
			(*bg.Members)[i].Player.(*BunkerPlayer).fill()
		}(iter)
	}

	wg.Wait()
}

func (bg *BunkerGame) openHiddenFeatures(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for i := range *bg.Members {
		if (*bg.Members)[i].Player.(*BunkerPlayer).isKicked {
			continue
		}

		memberIndex := i

		onSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			switch string(data) {
			case "1":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).biologicalParams.isHidden = false
			case "2":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).healthStatus.isHidden = false
			case "3":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).hobby.isHidden = false
			case "4":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).phobia.isHidden = false
			case "5":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).character.isHidden = false
			case "6":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).skill.isHidden = false
			case "7":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).knowledge.isHidden = false
			case "8":
				(*bg.Members)[memberIndex].Player.(*BunkerPlayer).baggage.isHidden = false
			default:
			}
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      (*bg.Members)[i].ChatID,
			Text:        text.ToOpen + (*bg.Members)[i].Player.Info(),
			ReplyMarkup: (*bg.Members)[i].Player.(*BunkerPlayer).Keyboard(b, onSelect),
		})
	}
}

func (bg *BunkerGame) sendAllInfo(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	data := []string{}

	for _, player := range *bg.Members {
		data = append(data, player.Name+player.Player.(*BunkerPlayer).View())
	}

	dataPool := sync.Pool{
		New: func() any {
			return make([]string, 0, MaximumMembersForBunker)
		},
	}

	for i := range *bg.Members {
		go func(current int) {
			defer wg.Done()

			dataToSend := dataPool.Get().([]string)
			defer func() {
				dataToSend = dataToSend[:0]
				dataPool.Put(&dataToSend)
			}()

			dataToSend = append(dataToSend, "Вы\n"+(*bg.Members)[current].Player.Info())
			dataToSend = append(dataToSend, data[:current]...)
			dataToSend = append(dataToSend, data[current+1:]...)

			opts := []paginator.Option{
				paginator.PerPage(1),
				paginator.WithCloseButton("Закрыть"),
			}

			p := paginator.New(b.(*bot.Bot), dataToSend, opts...)

			p.Show(ctx, b.(*bot.Bot), (*bg.Members)[current].ChatID)
		}(i)
	}
}

func (bg *BunkerGame) runVote(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	var innerWg sync.WaitGroup
	innerWg.Add(1)
	defer innerWg.Wait()

	for i := range *bg.Members {
		go func(current int) {
			onSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
				defer wg.Done()

				ChID, _ := strconv.Atoi(fmt.Sprintf("%d", data))
				player, ok := bg.Members.GetMember(int64(ChID))
				if !ok {
					logger.Log.Info("can't find member in lobby", zap.Int64("chat_id", mes.Message.Chat.ID), zap.Int64("mising_chat_id", int64(ChID)))
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: mes.Message.Chat.ID,
						Text: text.Convert(mes.Message.From.LanguageCode, text.SomethingWentWrong),
					})
					*bg.IsStarted = false
					bg.Members.SendAll(ctx, b, text.GameStoped)
					return
				}

				atomic.AddInt32(&player.Player.(*BunkerPlayer).votes, 1)

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   text.Convert(mes.Message.From.LanguageCode, text.VotedF, player.Name),
				})
			}

			kb := inline.New(b.(*bot.Bot))

			for _, member := range *bg.Members {
				if (*bg.Members)[current].ChatID == member.ChatID {
					continue
				}

				kb.Row().Button(member.Name, []byte(fmt.Sprintf("%d", member.ChatID)), onSelect)
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      (*bg.Members),
				Text:        text.ToKick,
				ReplyMarkup: kb,
			})
		}(i)
	}
}

func (bg *BunkerGame) twoMaxVotes() (User, User) {
	maxFirst := User{Player: &MafiaPlayer{votes: 0}}
	maxSecond := maxFirst

	for _, player := range *bg.Members {
		if maxFirst.Player.(*BunkerPlayer).votes < player.Player.(*BunkerPlayer).votes {
			maxSecond = maxFirst
			maxFirst = player
		} else if maxSecond.Player.(*BunkerPlayer).votes < player.Player.(*BunkerPlayer).votes {
			maxSecond = player
		}
	}

	return maxFirst, maxSecond
}

func (bg *BunkerGame) clearVotes() {
	for i := range *bg.Members {
		(*bg.Members)[i].Player.(*BunkerPlayer).votes = 0
	}
}

func (bg *BunkerGame) voteResults() string {
	result := ""

	for _, member := range *bg.Members {
		if len(result) == 0 {
			result += member.Name + ": " + fmt.Sprintf("%d", member.Player.(*BunkerPlayer).votes)
			continue
		}

		result += "\n" + member.Name + ": " + fmt.Sprintf("%d", member.Player.(*BunkerPlayer).votes)
	}

	return result
}

func (bg *BunkerGame) kick(u User) {
	index := bg.Members.FindMember(u)
	(*bg.Members)[index].Player.(*BunkerPlayer).isKicked = true
}

func (bg *BunkerGame) loadResults(repo *storage.Repository) {
	var wg sync.WaitGroup
	wg.Add(len(*bg.Members))

	for _, member := range *bg.Members {
		go func(user User) {
			defer wg.Done()
			err := repo.UpdateUserStatistic(user.ID, !user.Player.(*BunkerPlayer).isKicked)
			if err != nil {
                logger.Log.Error("can't update user statistic", zap.Int64("user_id", user.ID), zap.Error(err))
            }
		}(member)
	}

	wg.Wait()
}