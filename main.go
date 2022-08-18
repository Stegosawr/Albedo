package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/command"
	"github.com/stegosawr/Albedo/embed"
	"github.com/stegosawr/Albedo/reactions"
	"github.com/stegosawr/Albedo/static"
	"github.com/stegosawr/Albedo/utils"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	if _, err := dg.Channel(static.AnimeNewsChannelID); err == nil {
		timeTillInitialRun, _ := time.ParseDuration(fmt.Sprintf("%vh", 25-time.Now().Hour()))
		t := time.NewTimer(timeTillInitialRun)
		go func() {
			for {
				<-t.C
				t.Reset(24 * time.Hour)
				command.UpdateAnime(dg)
			}
		}()
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	regexList := map[string]*regexp.Regexp{
		// has to be done first -> unshortens twitter URLs so the bot can react
		static.TwitterProxy:  regexp.MustCompile(`https://t.co/\w+`),
		static.AmiAmi:        regexp.MustCompile(`https://www.amiami.(?:com|jp)/.+(?:[gs]code=[\w-]+|s_keywords=([\w-%]+))`),
		static.CuddlyOctopus: regexp.MustCompile(`https://cuddlyoctopus.com/product/[^/]+`),
		static.NHentai:       regexp.MustCompile(`^\+?[0-9]{5,6}`),
		static.SolarisJapan:  regexp.MustCompile(`https://solarisjapan\.com.*/products/.+`),
	}

	for k, v := range regexList {
		matchedRegex := v.FindString(m.Content)
		if len(matchedRegex) < 1 {
			continue
		}

		if k == static.TwitterProxy {
			newURL, err := utils.UnShortenURL(matchedRegex)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				break
			}

			m.Content = strings.ReplaceAll(m.Content, matchedRegex, newURL)
			continue
		}

		switch k {
		case static.AmiAmi, static.CuddlyOctopus, static.SolarisJapan:
			embeds, err := embed.Embed(s, m, k)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				break
			}
			for _, embed := range embeds {
				msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					break
				}

				if k == static.AmiAmi {
					s.MessageReactionAdd(m.ChannelID, msg.ID, "⬅️")
					s.MessageReactionAdd(m.ChannelID, msg.ID, "➡️")
				}

				if k == static.AmiAmi || k == static.CuddlyOctopus || k == static.SolarisJapan {
					s.MessageReactionAdd(m.ChannelID, msg.ID, "💶")
					s.MessageReactionAdd(m.ChannelID, msg.ID, "💴")
					s.MessageReactionAdd(m.ChannelID, msg.ID, "💵")
					s.MessageReactionAdd(m.ChannelID, msg.ID, "💷")
				}
			}
		case static.NHentai:
			command.NhentaiShow(s, m)
		}
		break
	}

	if !strings.HasPrefix(m.Content, static.Prefix) || m.Content == static.Prefix {
		return
	}

	switch strings.TrimPrefix(m.Content, static.Prefix) {
	case "help":
		command.Help(s, m)
	case "ping":
		command.Ping(s, m)
	case "dailyMedia":
		command.UpdateAnime(s)
	case "animeSites":
		s.ChannelMessageSend(m.ChannelID, "https://piracy.moe/")
	case "bestReleases":
		s.ChannelMessageSend(m.ChannelID, "https://releases.moe/")
	case "delmsg":
		command.DeleteAllMessagesInChannel(s, m.ChannelID)
	}
}

func messageReactionAdd(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {

	//ignore bot reactions
	if mra.UserID == s.State.User.ID {
		return
	}

	// pass data to the correct handlers
	err := reactions.Process(s, mra)
	if err != nil {
		s.ChannelMessageSend(mra.ChannelID, err.Error())
	}

}
