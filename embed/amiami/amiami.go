package amiami

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Stegosawr/apiapi"
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

const imgSite = "https://img.amiami.com"

var reFigureCode = regexp.MustCompile(`https://www\.amiami\.(?:com|jp)/.+([sg])code=([\w-]+)`)

type embeder struct{}

// New amiami embeder
func New() static.Embeder {
	return &embeder{}
}

// Embed from message content
func (e *embeder) Embed(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	matchedURL := reFigureCode.FindStringSubmatch(m.Content)
	if len(matchedURL) < 1 {
		return nil, errors.New("apiapi failed cannot parse CodeType and/or GCode")
	}

	codeType := apiapi.CodeTypeG
	if matchedURL[1] == "s" {
		codeType = apiapi.CodeTypeS
	}

	details, err := apiapi.GetItemByCode(codeType, matchedURL[2])
	if err != nil {
		return nil, err
	}

	if details.Item.CPriceTaxed == 0 {
		details.Item.CPriceTaxed = details.Item.Price
	}

	price := fmt.Sprintf("__%d__ JPY", details.Item.CPriceTaxed)
	if details.Item.CPriceTaxed > details.Item.Price {
		price = fmt.Sprintf("~~%d JPY~~-> __%d__ JPY | You save %d JPY", details.Item.CPriceTaxed, details.Item.Price, details.Item.CPriceTaxed-details.Item.Price)
	}

	status := "Release Date: " + details.Item.Releasedate
	if details.Item.Preorderitem == 1 {
		status = "*PRE-ORDER* - " + status
	}
	if details.Item.PreownAttention == 1 {
		status = "*PRE-OWNED* - " + status
	}

	return &discordgo.MessageEmbed{
		URL:         "https://www.amiami.com/eng/detail/?scode=" + details.Item.SCode,
		Title:       details.Item.SNameSimple,
		Description: details.Item.Spec,
		Image: &discordgo.MessageEmbedImage{
			URL: imgSite + details.Item.MainImageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("made by %s - Stock: %d - Images: 1/%d", details.Item.MakerName, details.Item.Stock, details.Item.ImageReviewnumber+1), // account for the main img
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
	}, nil
}

// GetNextImage to display of selected embed
func GetNextImage(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, static.ErrParsingMessageFailed.Error())
	}

	if len(msg.Embeds) != 1 {
		return
	}

	reImgNumb := regexp.MustCompile(`\d+$`)
	matchedNumbImages := reImgNumb.FindString(msg.Embeds[0].Footer.Text)

	// check if numb of images is supplied
	if matchedNumbImages == "" {
		return
	}

	numbImages, err := strconv.Atoi(matchedNumbImages)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, "get next image failed: "+err.Error())
	}

	reMainImgURL := regexp.MustCompile(`(\d+)/([^._]+)[^.]*\.(\w+)`)
	mainImgInfos := reMainImgURL.FindStringSubmatch(msg.Embeds[0].Image.URL) //1=category 2=FIGURE-CODE 3=ext
	if len(mainImgInfos) < 4 {
		return
	}

	imgURLs := []string{fmt.Sprintf("/images/product/main/%s/%s.%s", mainImgInfos[1], mainImgInfos[2], mainImgInfos[3])}
	for i := 1; numbImages > i; i++ {
		imgURLs = append(imgURLs, fmt.Sprintf("/images/product/review/%s/%s_%02d.%s", mainImgInfos[1], mainImgInfos[2], i, mainImgInfos[3]))
	}

	originalImgURL := strings.TrimPrefix(msg.Embeds[0].Image.URL, imgSite)
	originalImgURLIdx := 0
	for idx, i := range imgURLs {
		if i == originalImgURL {
			originalImgURLIdx = idx
		}
	}

	newImgIdx := originalImgURLIdx
	if mra.Emoji.Name == "⬅️" && originalImgURLIdx > 0 {
		newImgIdx = originalImgURLIdx - 1
	}
	if mra.Emoji.Name == "➡️" && originalImgURLIdx+1 <= numbImages-1 {
		newImgIdx = originalImgURLIdx + 1
		msg.Embeds[0].Image.URL = imgSite + imgURLs[originalImgURLIdx+1]
	}

	//fmt.Println(originalImgURL)
	//fmt.Println(msg.Embeds[0].Image.URL)
	s.MessageReactionRemove(mra.ChannelID, mra.MessageID, mra.Emoji.Name, mra.UserID)
	// no changes
	if newImgIdx == originalImgURLIdx {
		return
	}
	msg.Embeds[0].Image.URL = imgSite + imgURLs[newImgIdx]

	// update embed footer
	reCutFooter := regexp.MustCompile(`(.+) \d+/\d+$`)
	footerTemplate := reCutFooter.FindStringSubmatch(msg.Embeds[0].Footer.Text)
	if len(footerTemplate) == 2 {
		msg.Embeds[0].Footer.Text = fmt.Sprintf("%s %d/%d", footerTemplate[1], newImgIdx+1, numbImages)
	}

	// if img URL changed update the embed
	s.ChannelMessageEditEmbed(mra.ChannelID, mra.MessageID, msg.Embeds[0])
	if err != nil {
		fmt.Println(err)
	}

}
