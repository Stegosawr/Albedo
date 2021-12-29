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
		status = fmt.Sprintf("*%s* - Release Date: %s", strings.ToUpper(product.Product.Variants[0].Title), status)
	} else {
		status = fmt.Sprintf("*%s*", product.Product.Variants[0].Title)
	}

	description := ""
	for k, v := range product.Product.Info {
		if k == "Release Date" {
			continue
		}
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

// GetNextImage to display of selected embed
/*func GetNextImage(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, static.ErrParsingMessageFailed.Error())
	}

	if len(msg.Embeds) != 1 {
		return
	}

	reImgNumb := regexp.MustCompile(`(\d+)/(\d+)$`)
	matchedNumbImages := reImgNumb.FindString(msg.Embeds[0].Footer.Text)

	// check if numb of images is supplied
	if matchedNumbImages == "" {
		return
	}

	numbImages, err := strconv.Atoi(matchedNumbImages)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, "get next image failed: "+err.Error())
	}

	reImageFileName := regexp.MustCompile(`(\d+)\.\w+`)
	matchedImageFileName := reImageFileName.FindStringSubmatch(msg.Embeds[0].Image.URL) //1=filename -> is a number value
	if len(matchedImageFileName) < 2 {
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

}*/
