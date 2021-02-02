package command

import "github.com/bwmarrin/discordgo"

// Ping command implementation
func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
