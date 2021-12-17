package command

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	octopusapi "github.com/gan-of-culture/octopus-api"
)

var reProductURL = regexp.MustCompile(`https://cuddlyoctopus.com/product/[^/]+`)

func DakiShow(s *discordgo.Session, m *discordgo.MessageCreate) {
	matchedURL := reProductURL.FindString(m.Content)
	if matchedURL == "" {
		s.ChannelMessageSend(m.ChannelID, "invalid dakimakura URL")
		return
	}

	product, err := octopusapi.GetProductByURL(matchedURL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		URL:         product.URL,
		Title:       product.Name,
		Description: product.Description,
		Image: &discordgo.MessageEmbedImage{
			URL: product.MainImage,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("SKU: %d", product.Sku),
			IconURL: "https://cuddlyoctopus.com/wp-content/uploads/2016/03/cropped-Octodaki-00-Transparent-32x32.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Price:",
				Value: fmt.Sprintf("%s %s", product.Offers[0].Price, product.Offers[0].PriceCurrency),
			},
		},
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, msg.ID, "💶")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💴")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💵")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💷")
}
