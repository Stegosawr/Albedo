package arrow

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/embed/amiami"
	"github.com/stegosawr/Albedo/static"
)

var sites map[string]string

func init() {
	sites = map[string]string{
		"https://www.amiami.com/favicon.png": static.AmiAmi,
		"https://cuddlyoctopus.com/wp-content/uploads/2016/03/cropped-Octodaki-00-Transparent-32x32.png": static.CuddlyOctopus,
	}
}

type handler struct{}

// New ReactionHandler to handle arrow reactions
func New() static.ReactionHandler {
	return &handler{}
}

// Process arrow reactions
func (h *handler) Process(s *discordgo.Session, mra *discordgo.MessageReactionAdd) error {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		return static.ErrCurrencyConversionFailed
	}

	for _, embed := range msg.Embeds {
		if key, ok := sites[embed.Footer.IconURL]; ok {
			switch key {
			case static.AmiAmi:
				amiami.GetNextImage(s, mra)
			}
		}
	}
	return nil
}
