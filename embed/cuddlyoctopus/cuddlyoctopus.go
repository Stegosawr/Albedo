package cuddlyoctopus

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	octopusapi "github.com/gan-of-culture/octopus-api"
	"github.com/stegosawr/Albedo/static"
)

var reProductURL = regexp.MustCompile(`https://cuddlyoctopus.com/product/[^/]+`)

type embeder struct{}

// New cuddlyoctopus embeder
func New() static.Embeder {
	return &embeder{}
}

// Embed from message content
func (e *embeder) Embed(s *discordgo.Session, m *discordgo.MessageCreate) ([]*discordgo.MessageEmbed, error) {
	matchedURL := reProductURL.FindString(m.Content)
	if matchedURL == "" {
		return nil, static.ErrURLParseFailed
	}

	product, err := octopusapi.GetProductByURL(matchedURL)
	if err != nil {
		return nil, err
	}

	imageURL := product.NSFWImage
	if imageURL == "" {
		imageURL = product.MainImage
	}

	if strings.Contains(m.Content, "#sfw") {
		imageURL = product.MainImage
	}

	return []*discordgo.MessageEmbed{
		{
			URL:         product.URL,
			Title:       product.Name,
			Description: product.Description,
			Image: &discordgo.MessageEmbedImage{
				URL: imageURL,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
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
		},
	}, nil
}
