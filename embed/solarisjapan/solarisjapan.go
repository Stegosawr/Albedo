package solarisjapan

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Stegosawr/solarisapi"
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

var reProductURL = regexp.MustCompile(`https://solarisjapan\.com.*/products/.+`)

type embeder struct{}

// New solarisjapan embeder
func New() static.Embeder {
	return &embeder{}
}

// Embed from message content
func (e *embeder) Embed(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	matchedURL := reProductURL.FindString(m.Content)
	if matchedURL == "" {
		return nil, static.ErrURLParseFailed
	}

	product, err := solarisapi.GetItemByURL(matchedURL)
	if err != nil {
		return nil, err
	}

	status, ok := product.Product.Info["Release Date"]
	if ok {
		status = fmt.Sprintf("*%s* - Release Date: %s", product.Product.Variants[0].Title, status)
	} else {
		status = fmt.Sprintf("*%s*", product.Product.Variants[0].Title)
	}

	description := ""
	for k, v := range product.Product.Info {
		description = fmt.Sprintf("%s\n%s: %s", description, k, v)
	}
	description = strings.TrimSpace(description)

	return &discordgo.MessageEmbed{
		URL:         matchedURL,
		Title:       product.Product.Title,
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: product.Product.Image.Src,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("made by: %s - Images: 1/%d", product.Product.Vendor, len(product.Product.Images)),
			IconURL: "https://cdn.shopify.com/s/files/1/0318/2649/t/54/assets/favicon.ico?v=8178696738005321811",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Price:",
				Value: fmt.Sprintf("%s JPY", product.Product.Variants[0].Price),
			},
			{
				Name:  "Status:",
				Value: status,
			},
		},
	}, nil
}
