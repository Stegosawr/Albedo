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
	"github.com/stegosawr/Albedo/static"
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
				select {
				case <-t.C:
					t.Reset(24 * time.Hour)
					command.UpdateAnime(dg)
				}
			}
		}()
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
		"amiami":  regexp.MustCompile(`https://www.amiami.com/.+[gs]code=[\w-]+`),
		"nhentai": regexp.MustCompile(`^\+?[0-9]{5,6}`),
	}

	for k, v := range regexList {
		matchedRegex := v.FindString(m.Content)
		if len(matchedRegex) < 1 {
			continue
		}

		switch k {
		case "amiami":
			command.FigureShow(s, m)
		case "nhentai":
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
	case "watch anime":
		s.ChannelMessageSend(m.ChannelID, "https://piracy.moe/")
	case "delmsg":
		command.DeleteAllMessagesInChannel(s, m.ChannelID)
	}
}

func messageReactionAdd(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {

	//ignore bot reactions
	if mra.UserID == s.State.User.ID {
		return
	}

	//only do somthing if it matches these emotes
	switch mra.Emoji.Name {
	case "💵", "💴", "💶", "💷":
		command.ConvertCurrencies(s, mra)
	case "⬅️", "➡️":
		command.GetNextImage(s, mra)
		return
	default:
		return
	}

}
