package embed

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/embed/amiami"
	"github.com/stegosawr/Albedo/embed/cuddlyoctopus"
	"github.com/stegosawr/Albedo/static"
)

var embedersMap map[string]static.Embeder

func init() {
	embedersMap = map[string]static.Embeder{
		static.AmiAmi:        amiami.New(),
		static.CuddlyOctopus: cuddlyoctopus.New(),
	}
}

// Embed call the other Embeders
func Embed(s *discordgo.Session, m *discordgo.MessageCreate, key string) (*discordgo.MessageEmbed, error) {
	embeder := embedersMap[key]
	if embeder == nil {
		return nil, nil
	}
	return embeder.Embed(s, m)
}
