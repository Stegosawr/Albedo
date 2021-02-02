package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

// NhentaiShow enter a magical number to get some magical information
func NhentaiShow(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s/", "https://nhentai.net/g/", strings.TrimPrefix(m.Content, static.Prefix)))
}
