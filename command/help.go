package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

// Help show information about all commands
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpEmbedField := []*discordgo.MessageEmbedField{}
	for key, command := range static.Commands {
		helpEmbedField = append(helpEmbedField, &discordgo.MessageEmbedField{
			Name:  key,
			Value: command.Description,
		})
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  "Command list:",
		Fields: helpEmbedField,
	})
}
