package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Stegosawr/apiapi"
	"github.com/bwmarrin/discordgo"
)

const site = "https://amiami.com"
const imgSite = "https://img.amiami.com"

var reFigureCode = regexp.MustCompile(`https://www\.amiami\.com/.+([sg])code=([\w-]+)`)
var reCurrencies = regexp.MustCompile(`([\d,.]+)(__|~~)? ([A-Z]{3})`)

var currencies = map[string]string{
	"💵": "USD",
	"💴": "JPY",
	"💶": "EUR",
	"💷": "GBP",
}

var exchangeRates apiapi.CurrencyLayer

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

	imageURLsAsCSV := details.Item.MainImageURL
	for _, reviewImg := range details.Embedded.ReviewImages {
		imageURLsAsCSV = imageURLsAsCSV + "," + reviewImg.ImageURL
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		URL:         "https://www.amiami.com/eng/detail/?scode=" + details.Item.SCode,
		Title:       details.Item.SNameSimple,
		Description: details.Item.Spec,
		Image: &discordgo.MessageEmbedImage{
			URL: imgSite + details.Item.MainImageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("made by %s - Stock: %d - Images: 1/%d", details.Item.MakerName, details.Item.Stock, details.Item.ImageReviewnumber),
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
	if err != nil {
		fmt.Println(err)
	}

	s.MessageReactionAdd(m.ChannelID, msg.ID, "⬅️")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "➡️")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💶")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💴")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💵")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "💷")
}

// ConvertCurrencies in Message
func ConvertCurrencies(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, "currency conversion failed: "+err.Error())
	}

	exchangeRates, err = apiapi.GetCurrencyLayer()
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, "currency conversion failed: "+err.Error())
	}

	// remove the users reaction
	s.MessageReactionRemove(mra.ChannelID, mra.MessageID, mra.Emoji.Name, mra.UserID)

	contentToConvert := reCurrencies.FindAllStringSubmatch(msg.Content, -1)
	if len(contentToConvert) > 0 {
		newContent := convertCurr(msg.Content, contentToConvert[0][3], mra.Emoji.Name)
		if msg.Author.ID == s.State.User.ID {
			_, err = s.ChannelMessageEdit(mra.ChannelID, mra.MessageID, newContent)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			s.ChannelMessageSend(mra.ChannelID, newContent)
		}
	}

	if len(msg.Embeds) == 0 {
		return
	}

	for _, e := range msg.Embeds {
		contentToConvert = reCurrencies.FindAllStringSubmatch(e.Description, -1)
		if len(contentToConvert) > 0 {
			e.Description = convertCurr(e.Description, contentToConvert[0][3], mra.Emoji.Name)
		}

		contentToConvert = reCurrencies.FindAllStringSubmatch(e.Footer.Text, -1)
		if len(contentToConvert) > 0 {
			e.Footer.Text = convertCurr(e.Footer.Text, contentToConvert[0][3], mra.Emoji.Name)
		}

		for _, f := range e.Fields {
			contentToConvert = reCurrencies.FindAllStringSubmatch(f.Value, -1)
			if len(contentToConvert) > 0 {
				f.Value = convertCurr(f.Value, contentToConvert[0][3], mra.Emoji.Name)
			}
		}

		if msg.Author.ID == s.State.User.ID {
			_, err = s.ChannelMessageEditEmbed(mra.ChannelID, mra.MessageID, e)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			s.ChannelMessageSendEmbed(mra.ChannelID, e)
		}
	}

}

func convertCurr(content, source, target string) string {
	toConvert := reCurrencies.FindAllStringSubmatch(content, -1)

	var amount float64
	var err error
	for _, t := range toConvert {
		amount, err = strconv.ParseFloat(t[1], 64)
		if err != nil {
			fmt.Println("curreny conversion error can't parse amount from msg content")
			continue
		}

		if source != "USD" {
			amount = amount / getExchangeRate(source)
		}
		//fmt.Println(amount)
		if target != "💵" {
			amount = amount * getExchangeRate(currencies[target])
		}
		//fmt.Println(amount)

		// now amount is in USD
		if target == "💴" {
			content = strings.ReplaceAll(content, t[0], fmt.Sprintf("%.0f%s %s", amount, t[2], currencies[target]))
		} else {
			content = strings.ReplaceAll(content, t[0], fmt.Sprintf("%.2f%s %s", amount, t[2], currencies[target]))
		}

		//fmt.Println(content)
	}

	return content
}

func getExchangeRate(currkey string) float64 {
	switch currkey {
	case "EUR":
		return exchangeRates.Quotes.USDEUR
	case "JPY":
		return exchangeRates.Quotes.USDJPY
	case "GBP":
		return exchangeRates.Quotes.USDGBP
	}
	return 0.0
}

// GetNextImage to display of selected embed
func GetNextImage(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, "get next image failed: "+err.Error())
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

	imgURLs := []string{strings.TrimPrefix(msg.Embeds[0].Image.URL, imgSite)}
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
