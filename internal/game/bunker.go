package game

import (
	"context"
	"gffbot/internal/text"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/slider"
)

type BunkerPlayerFeature struct {
	info     string
	isHidden bool
}

func (bf *BunkerPlayerFeature) toString() string {
	if bf.isHidden {
		return "[hidden] " + bf.info
	} else {
		return "[shown]  " + bf.info
	}
}

func (bf *BunkerPlayerFeature) view() string {
	if bf.isHidden {
		return "hidden"
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
	skill            BunkerPlayerFeature
	knowledge        BunkerPlayerFeature
	baggage          BunkerPlayerFeature

	actionCard    any
	conditionCard any

	votes    int
	isKicked bool
}

func (bp *BunkerPlayer) GetInfo() string {
	return text.GetConvertToLang(*bp.Lang, text.Profession) + ":\n" +
		bp.profession +
		"\n" + text.GetConvertToLang(*bp.Lang, text.BoilogicalParams) + ":\n" +
		bp.biologicalParams.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.HealthStatus) + ":\n" +
		bp.healthStatus.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Hobby) + ":\n" +
		bp.hobby.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Phopia) + ":\n" +
		bp.phobia.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Character) + ":\n" +
		bp.character.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Skill) + ":\n" +
		bp.skill.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Knowledge) + ":\n" +
		bp.knowledge.toString() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Baggage) + ":\n" +
		bp.baggage.toString() +
		"\n\n" + text.GetConvertToLang(*bp.Lang, text.ActionCard) + ":\n" +
		"TODO" +
		"\n" + text.GetConvertToLang(*bp.Lang, text.ConditionCard) + ":\n" +
		"TODO"
}

func (bp *BunkerPlayer) GetView() string {
	return text.GetConvertToLang(*bp.Lang, text.Profession) + ":\n" +
		bp.profession +
		"\n" + text.GetConvertToLang(*bp.Lang, text.BoilogicalParams) + ":\n" +
		bp.biologicalParams.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.HealthStatus) + ":\n" +
		bp.healthStatus.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Hobby) + ":\n" +
		bp.hobby.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Phopia) + ":\n" +
		bp.phobia.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Character) + ":\n" +
		bp.character.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Skill) + ":\n" +
		bp.skill.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Knowledge) + ":\n" +
		bp.knowledge.view() +
		"\n" + text.GetConvertToLang(*bp.Lang, text.Baggage) + ":\n" +
		bp.baggage.view()
}

func (bp *BunkerPlayer) getKeyboard(b Bot, onSelect inline.OnSelect) *inline.Keyboard {
	kb := inline.New(b.(*bot.Bot))

	if bp.biologicalParams.isHidden {
		kb.Row().Button("Boilogical params", []byte("1"), onSelect)
	}
	if bp.healthStatus.isHidden {
		kb.Row().Button("Helth status", []byte("2"), onSelect)
	}
	if bp.hobby.isHidden {
		kb.Row().Button("Hobby", []byte("3"), onSelect)
	}
	if bp.phobia.isHidden {
		kb.Row().Button("Phobia", []byte("4"), onSelect)
	}
	if bp.character.isHidden {
		kb.Row().Button("Character", []byte("5"), onSelect)
	}
	if bp.skill.isHidden {
		kb.Row().Button("Skill", []byte("6"), onSelect)
	}
	if bp.knowledge.isHidden {
		kb.Row().Button("Knowledge", []byte("7"), onSelect)
	}
	if bp.baggage.isHidden {
		kb.Row().Button("Baggage", []byte("8"), onSelect)
	}

	return kb
}

type BunkerGame struct {
	disastre  string
	IsStarted *bool
	Members   *Users
}

func (bg *BunkerGame) StartGame(ctx context.Context, b Bot) {
	// fill roles and disastre

	var wg sync.WaitGroup

	countOfMembers := len(*bg.Members)

	for *bg.IsStarted {
		// open hidden features
		bg.openHiddenFeatures(ctx, b, &wg)

		wg.Add(countOfMembers)

		bg.sendSliders(ctx, b, &wg)

		wg.Wait()

		break
	}

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
			Text:        (*bg.Members)[i].Player.GetInfo(),
			ReplyMarkup: (*bg.Members)[i].Player.(*BunkerPlayer).getKeyboard(b, onSelect),
		})
	}
}

func (bg *BunkerGame) sendSliders(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for _, member := range *bg.Members {
		if member.Player.(*BunkerPlayer).isKicked {
			continue
		}

		slides := []slider.Slide{{Text: member.Player.GetInfo()}}
		membersList := []int64{}

		for _, other := range *bg.Members {
			if member.ChatID == other.ChatID || other.Player.(*BunkerPlayer).isKicked {
				continue
			}

			slides = append(slides, slider.Slide{Text: other.Name + "\n" + other.Player.(*BunkerPlayer).GetView()})
			membersList = append(membersList, other.ChatID)
		}

		onVoteSelect := func(ctx context.Context, b *bot.Bot, message models.MaybeInaccessibleMessage, item int) {
			defer wg.Done()

			// TODO
			if membersList[item] == message.Message.Chat.ID {
				// message: can't vote for yourself!
				return
			}

			(*bg.Members)[bg.Members.FindMember(
				User{ChatID: message.Message.Chat.ID},
			)].Player.(*BunkerPlayer).votes++

			// TODO
			// message: voted for %s
		}

		opts := []slider.Option{
			slider.OnSelect("Vote", true, onVoteSelect),
		}

		sl := slider.New(b.(*bot.Bot), slides, opts...)

		sl.Show(ctx, b.(*bot.Bot), member.ChatID)
	}
}
