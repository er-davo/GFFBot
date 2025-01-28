package game

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/slider"
)

type BunkerPlayerFeature struct {
	info		string
	isHidden	bool
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
	lang				*string

	profession			string

	biologicalParams	BunkerPlayerFeature
	healthStatus		BunkerPlayerFeature
	hobby				BunkerPlayerFeature
	phobia				BunkerPlayerFeature
	character			BunkerPlayerFeature
	skill				BunkerPlayerFeature
	knowledge			BunkerPlayerFeature
	baggage				BunkerPlayerFeature

	actionCard			any
	conditionCard		any

	votes				int
	isKicked			bool
}

func (bp *BunkerPlayer) GetInfo() string {
	return "Profession:\n" +
			bp.profession +
			"\nBoilogical params:\n" +
			bp.biologicalParams.toString() +
			"\nHealth status:\n" +
			bp.healthStatus.toString() +
			"\nHobby:\n" +
			bp.hobby.toString() +
			"\nPhopia:\n" +
			bp.phobia.toString() +
			"\nCharacter:\n" +
			bp.character.toString() +
			"\nSkill:\n" +
			bp.skill.toString() +
			"\nKnowledge:\n" +
			bp.knowledge.toString() +
			"\nBaggage:\n" +
			bp.baggage.toString() +
			"\n\nAction card:\n" +
			"TODO" +
			"\nCondition card:\n" +
			"TODO"
}

func (bp *BunkerPlayer) GetView() string {
	return "Profession:\n" +
			bp.profession +
			"\nBoilogical params:\n" +
			bp.biologicalParams.view() +
			"\nHealth status:\n" +
			bp.healthStatus.view() +
			"\nHobby:\n" +
			bp.hobby.view() +
			"\nPhopia:\n" +
			bp.phobia.view() +
			"\nCharacter:\n" +
			bp.character.view() +
			"\nSkill:\n" +
			bp.skill.view() +
			"\nKnowledge:\n" +
			bp.knowledge.view() +
			"\nBaggage:\n" +
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
	disastre	string
	isStarted	*bool
    members     *Users
}

func (bg *BunkerGame) startGame(ctx context.Context, b Bot) {
	// fill roles and disastre

	var wg sync.WaitGroup

	countOfMembers := len(*bg.members)

	for *bg.isStarted {
		// open hidden features
		bg.openHiddenFeatures(ctx, b, &wg)

		wg.Add(countOfMembers)

		bg.sendSliders(ctx, b, &wg)

		wg.Wait()

		break
	}

}

func (bg *BunkerGame) openHiddenFeatures(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for i := range *bg.members {
		if (*bg.members)[i].player.(*BunkerPlayer).isKicked {
			continue
		}
		
		memberIndex := i

		onSelect := func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			defer wg.Done()

			switch string(data) {
			case "1":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).biologicalParams.isHidden = false
			case "2":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).healthStatus.isHidden = false
			case "3":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).hobby.isHidden = false
			case "4":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).phobia.isHidden = false
			case "5":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).character.isHidden = false
			case "6":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).skill.isHidden = false
			case "7":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).knowledge.isHidden = false
			case "8":
				(*bg.members)[memberIndex].player.(*BunkerPlayer).baggage.isHidden = false
			default:
			}
		}
		
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:			(*bg.members)[i].ChatID,
			Text:			(*bg.members)[i].player.GetInfo(),
			ReplyMarkup:	(*bg.members)[i].player.(*BunkerPlayer).getKeyboard(b, onSelect),
		})
	}
}

func (bg *BunkerGame) sendSliders(ctx context.Context, b Bot, wg *sync.WaitGroup) {
	for _, member := range *bg.members {
		if member.player.(*BunkerPlayer).isKicked {
			continue
		}

		slides := []slider.Slide{{Text: member.player.GetInfo()}}
		membersList := []int64{}

		for _, other := range *bg.members {
			if member.ChatID == other.ChatID || other.player.(*BunkerPlayer).isKicked {
				continue
			}

			slides = append(slides, slider.Slide{Text: other.Name + "\n" + other.player.(*BunkerPlayer).GetView()})
			membersList = append(membersList, other.ChatID)
		}

		onVoteSelect := func(ctx context.Context, b *bot.Bot, message models.MaybeInaccessibleMessage, item int)  {
			defer wg.Done()

			// TODO
			if membersList[item] == message.Message.Chat.ID {
				// message: can't vote for yourself!
				return
			}
			
			(*bg.members)[bg.members.findMember(
				User{ChatID: message.Message.Chat.ID},
			)].player.(*BunkerPlayer).votes++

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