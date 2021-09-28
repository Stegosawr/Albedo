package command

import (
	"fmt"
	"regexp"

	"github.com/Stegosawr/apiapi"
	"github.com/bwmarrin/discordgo"
)

const site = "https://amiami.com"
const imgSite = "https://img.amiami.com/"

var reFigureCode = regexp.MustCompile(`https://www\.amiami\.com/.+([sg])code=([\w-]+)`)

func FigureShow(s *discordgo.Session, m *discordgo.MessageCreate) {

	matchedURL := reFigureCode.FindStringSubmatch(m.Content)
	if len(matchedURL) < 1 {
		s.ChannelMessageSend(m.ChannelID, "apiapi failed cannot parse CodeType and/or GCode")
	}

	codeType := apiapi.CodeTypeG
	if matchedURL[1] == "s" {
		codeType = apiapi.CodeTypeS
	}

	details, err := apiapi.GetItemByCode(codeType, matchedURL[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "apiapi failed: "+err.Error())
	}

	price := fmt.Sprintf("__%d__", details.Item.CPriceTaxed)
	if details.Item.CPriceTaxed > details.Item.Price {
		price = fmt.Sprintf("~~%d~~ -> __%d__ JPY | You save %d JPY", details.Item.CPriceTaxed, details.Item.Price, details.Item.CPriceTaxed-details.Item.Price)
	}

	status := "Release Date: " + details.Item.Releasedate
	if details.Item.Preorderitem == 1 {
		status = "*PRE-ORDER* - " + status
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		URL:         "https://www.amiami.com/eng/detail/?scode=" + details.Item.SCode,
		Title:       details.Item.SNameSimple,
		Description: details.Item.Spec,
		Image: &discordgo.MessageEmbedImage{
			URL: imgSite + details.Item.MainImageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("made by %s - Stock: %d", details.Item.MakerName, details.Item.Stock),
			IconURL: "https://www.amiami.com/favicon.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Price:",
				Value: price,
			},
			{
				Name:  "Status:",
				Value: status,
			},
		},
	})
}
