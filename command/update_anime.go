package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/request"
	"github.com/stegosawr/Albedo/static"
)

// AnimeCard = 1 Anime
type AnimeCard struct {
	URL        string
	Title      string
	ThumbURL   string
	Episode    string
	EpisodeURL string
}

// AnimeCollection Grouped by day
type AnimeCollection struct {
	Date    string
	IdxFrom int
	IdxTo   int
	Animes  []AnimeCard
}

// AnimeSchedule = Days Grouped
type AnimeSchedule struct {
	Days []AnimeCollection
}

// UpdateAnime Anime Schedule
func UpdateAnime(s *discordgo.Session) {

	_ = DeleteAllMessagesInChannel(s, static.AnimeNewsChannelID)

	htmlBody, err := request.Get("https://anidb.net/anime/schedule")
	if err != nil {
		s.ChannelMessageSend(static.AnimeNewsChannelID, "Can not reach https://anidb.net/anime/schedule")
	}

	re := regexp.MustCompile("\"(/anime/[0-9]{1,5})")
	matchedURLs := re.FindAllStringSubmatch(htmlBody, -1)

	re = regexp.MustCompile("aid=[0-9]*\">(.*) -? ?([a-zA-Z0-9]*)")
	matchedTitles := re.FindAllStringSubmatch(htmlBody, -1)

	re = regexp.MustCompile("https://cdn-eu.anidb.net/images/[^/]*/[^-]*")
	ThumbURLs := re.FindAllString(htmlBody, -1)

	re = regexp.MustCompile("/episode/[0-9]*/\\?aid=[0-9]*")
	EpisodeURLs := re.FindAllString(htmlBody, -1)

	AnimeCards := []AnimeCard{}
	ep := ""
	for idx, matchedURL := range matchedURLs {
		if matchedURL[1] == "" {
			continue
		}
		ep = ""
		if len(matchedTitles[idx]) == 3 {
			ep = matchedTitles[idx][2]
		}

		AnimeCards = append(AnimeCards, AnimeCard{
			URL:        matchedURL[1],
			Title:      strings.TrimRight(matchedTitles[idx][1], " -"),
			ThumbURL:   ThumbURLs[idx],
			Episode:    ep,
			EpisodeURL: EpisodeURLs[idx],
		})
	}

	re = regexp.MustCompile("[0-9]*-[0-9]*-[0-9]*, [^<]*")
	dates := re.FindAllString(htmlBody, -1)

	animeSchedule := AnimeSchedule{}
	datePosition := 0
	for idx, date := range dates {
		datePosition = strings.Index(htmlBody, date)
		if datePosition == -1 {
			s.ChannelMessageSend(static.AnimeNewsChannelID, "Can't group by Date - internal error")
			return
		}

		animeCollection := AnimeCollection{
			Date:    date,
			IdxFrom: datePosition,
		}

		if idx == 0 {
			animeCollection.IdxFrom = 0
		}

		if idx == len(dates)-1 {
			animeCollection.IdxTo = len(htmlBody)
		} else {
			animeCollection.IdxTo = strings.Index(htmlBody, dates[idx+1])
		}

		//fmt.Printf("idx: %d; from: %d; to: %d;\n", idx, animeCollection.IdxFrom, animeCollection.IdxTo)

		animeSchedule.Days = append(animeSchedule.Days, animeCollection)
	}

	if static.AnimeOfCurrentDay {
		for _, day := range animeSchedule.Days {
			if strings.Contains(day.Date, "today") {
				animeSchedule.Days = []AnimeCollection{day}
				break
			}
		}
	}

	episodeURLPosition := 0
	for _, anime := range AnimeCards {
		episodeURLPosition = strings.LastIndex(htmlBody, anime.URL)
		if episodeURLPosition == -1 {
			s.ChannelMessageSend(static.AnimeNewsChannelID, "Can't find url position - internal error")
			return
		}

		for idx, collection := range animeSchedule.Days {
			if episodeURLPosition > collection.IdxFrom && episodeURLPosition < collection.IdxTo {
				//fmt.Printf("title: %s; epURL: %s; episodeUrlPos: %d; from: %d; to: %d;\n", anime.Title, anime.EpisodeURL, episodeURLPosition, collection.IdxFrom, collection.IdxTo)
				animeSchedule.Days[idx].Animes = append(animeSchedule.Days[idx].Animes, anime)
			}
		}

	}
	for _, day := range animeSchedule.Days {
		for _, a := range day.Animes {
			_, err = s.ChannelMessageSendEmbed(static.AnimeNewsChannelID, &discordgo.MessageEmbed{
				URL:         fmt.Sprintf("%s%s", "https://anidb.net", a.EpisodeURL),
				Title:       a.Title,
				Description: fmt.Sprintf("Episode %s", a.Episode),
				Image: &discordgo.MessageEmbedImage{
					URL: a.ThumbURL,
				},
			})
			if err != nil {
				fmt.Println(err)
			}
		}
		break
	}

}
