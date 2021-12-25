package reactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/reactions/arrow"
	"github.com/stegosawr/Albedo/reactions/currency"
	"github.com/stegosawr/Albedo/static"
)

var handlersMap map[string]static.ReactionHandler

func init() {
	currencyHandler := currency.New()
	arrowHandler := arrow.New()

	handlersMap = map[string]static.ReactionHandler{
		"💵":  currencyHandler,
		"💴":  currencyHandler,
		"💶":  currencyHandler,
		"💷":  currencyHandler,
		"⬅️": arrowHandler,
		"➡️": arrowHandler,
	}
}

// Process generally and then call the other implementaions
func Process(s *discordgo.Session, mra *discordgo.MessageReactionAdd) error {
	handler := handlersMap[mra.Emoji.Name]
	if handler == nil {
		return nil
	}
	return handler.Process(s, mra)
}
