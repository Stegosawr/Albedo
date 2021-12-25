package static

import "github.com/bwmarrin/discordgo"

// Command Describes a command
type Command struct {
	Description        string
	SpecialPermissions string
}

// Embeder of different types of content mostly website content
type Embeder interface {
	Embed(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error)
}

// ReactionHandler of reactions to messages
type ReactionHandler interface {
	Process(s *discordgo.Session, mra *discordgo.MessageReactionAdd) error
}
