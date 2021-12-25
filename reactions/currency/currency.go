package currency

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Stegosawr/currency"
	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

var exchangeRates map[string]float64
var currencies = map[string]string{
	"💵": "USD",
	"💴": "JPY",
	"💶": "EUR",
	"💷": "GBP",
}
var reCurrencies = regexp.MustCompile(`([\d,.]+)(__|~~)? ([A-Z]{3})`)

type handler struct{}

// New ReactionHandler to handle currency reactions
func New() static.ReactionHandler {
	return &handler{}
}

// Process currency reactions
func (h *handler) Process(s *discordgo.Session, mra *discordgo.MessageReactionAdd) error {
	msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
	if err != nil {
		return static.ErrParsingMessageFailed
	}

	exchangeRates, err = currency.GetCurrencyRates()
	if err != nil {
		return err
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
		return nil
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
	return nil
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
			amount = amount * getExchangeRate(source)
		}
		//fmt.Println(amount)
		if target != "💵" {
			amount = amount / getExchangeRate(currencies[target])
		}
		//fmt.Println(amount)

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
	if rate, ok := exchangeRates[currkey]; ok {
		return rate
	}
	return 0.0
}
