package command

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// DeleteMessage from channel
func DeleteMessage(s *discordgo.Session, m *discordgo.Message) error {
	return s.ChannelMessageDelete(m.ChannelID, m.ID)
}

// DeleteMessages from channel
func DeleteMessages(s *discordgo.Session, m []*discordgo.Message) []error {
	fmt.Println(len(m))
	errors := []error{}
	for _, msg := range m {
		time.Sleep(1 * time.Second)
		err := DeleteMessage(s, msg)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// DeleteAllMessagesInChannel (risky)
func DeleteAllMessagesInChannel(s *discordgo.Session, ChannelID string) error {
	for {
		if c, _ := s.ChannelMessages(ChannelID, 0, "", "", ""); len(c) > 0 {
			errors := DeleteMessages(s, c)
			if len(errors) > 0 {
				return errors[0]
			}
			continue
		}
		break
	}
	return nil
}
