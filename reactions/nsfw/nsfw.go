package nsfw

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

type handler struct{}

// New ReactionHandler to handle nsfw reactions
func New() static.ReactionHandler {
	return &handler{}
}

// Process nsfw reactions
func (h *handler) Process(s *discordgo.Session, mra *discordgo.MessageReactionAdd) error {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		return static.ErrCurrencyConversionFailed
	}

	for _, embed := range msg.Embeds {
		fmt.Println(embed.Thumbnail.URL)
	}
	return nil
}
