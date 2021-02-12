package command

import (
	"fmt"
	"regexp"
	"strings"
	"time"

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
	Date       string
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

	_ = DeleteAllMessagesInChannel(s, static.HentaiNewsChannelID)

	animeSchedule, err := scraperAnime()
	if err != nil {
		s.ChannelMessageSend(static.AnimeNewsChannelID, err.Error())
		return
	}
	hentaiSchdule, err := scraperHentai()
	if err != nil {
		s.ChannelMessageSend(static.AnimeNewsChannelID, err.Error())
		return
	}

	animeSchedules := []AnimeSchedule{animeSchedule, hentaiSchdule}
	for idx, schedule := range animeSchedules {
		channel := static.AnimeNewsChannelID
		switch idx {
		case 0:
			s.ChannelMessageSend(channel, "--------------------------\n  Anime of Today\n--------------------------")
		case 1:
			channel = static.HentaiNewsChannelID
			s.ChannelMessageSend(channel, "--------------------------\n  Hentai of Today\n--------------------------")
			if len(animeSchedules[idx].Days) == 0 {
				s.ChannelMessageSend(channel, fmt.Sprintf("No releases today!. Check release schedule here %s", fmt.Sprintf("https://www.underhentai.net/releases-%d/", time.Now().Year())))
			}
		}
		for _, day := range schedule.Days {
			for _, a := range day.Animes {
				time.Sleep(1 * time.Second)
				_, err = s.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
					URL:         a.URL,
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

}

func scraperAnime() (AnimeSchedule, error) {

	htmlBody, err := request.Get("https://anidb.net/anime/schedule")
	if err != nil {
		return AnimeSchedule{}, fmt.Errorf("%s %s", static.AnimeNewsChannelID, "Can not reach https://anidb.net/anime/schedule")
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
			return AnimeSchedule{}, fmt.Errorf("%s %s", static.AnimeNewsChannelID, "Can't group by Date - internal error")
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
			return AnimeSchedule{}, fmt.Errorf("%s %s", static.AnimeNewsChannelID, "Can't find url position - internal error")
		}

		for idx, collection := range animeSchedule.Days {
			if episodeURLPosition > collection.IdxFrom && episodeURLPosition < collection.IdxTo {
				anime.URL = fmt.Sprintf("%s%s", "https://anidb.net", EpisodeURLs[idx])
				animeSchedule.Days[idx].Animes = append(animeSchedule.Days[idx].Animes, anime)
			}
		}

	}
	return animeSchedule, nil
}

func scraperHentai() (AnimeSchedule, error) {
	now := time.Now()

	htmlBody, err := request.Get(fmt.Sprintf("https://www.underhentai.net/releases-%d/", now.Year()))
	if err != nil {
		return AnimeSchedule{}, fmt.Errorf("%s %s", static.AnimeNewsChannelID, "Can not reach https://anidb.net/anime/schedule")
	}

	re := regexp.MustCompile("\"article-section[^<]*\\s*<[^\"]*\"([^\"]*)\"\\s*[^\"]*\"([^\"]*)[^/]*([^\"]*)")
	matchedImgTag := re.FindAllStringSubmatch(htmlBody, -1) // 1=url 2=title 3=imgUrl

	re = regexp.MustCompile("article-footer>[^:]*:\\s?([0-9]*)[^:]*:\\s?([^<]*)")
	matchedFooters := re.FindAllStringSubmatch(htmlBody, -1) // 1=epNo 2=releaseDate

	animes := []AnimeCard{}
	for idx, matchedImgTag := range matchedImgTag {
		animes = append(animes, AnimeCard{
			URL:      fmt.Sprintf("%s%s", "https://www.underhentai.net", matchedImgTag[1]),
			Title:    matchedImgTag[2],
			ThumbURL: fmt.Sprintf("%s%s", "https:", matchedImgTag[3]),
			Episode:  matchedFooters[idx][1],
			Date:     strings.ReplaceAll(matchedFooters[idx][2], "/", "."),
		})
	}

	dates := []string{}
	for key := range matchedFooters {
		dates = append(dates, matchedFooters[key][2])
	}

	dates = removeDuplicatesUnordered(dates)

	animeCollection := []AnimeCollection{}
	for _, date := range dates {
		foramttedDate := strings.ReplaceAll(date, "/", ".")
		realDate := fmt.Sprintf("%s.%s.%s", addZeroTo2DigitNum(now.Day()), addZeroTo2DigitNum(int(now.Month())), addZeroTo2DigitNum(now.Year()))
		if realDate != foramttedDate && static.AnimeOfCurrentDay {
			continue
		}
		animeCollection = append(animeCollection, AnimeCollection{
			Date: foramttedDate,
		})
	}

	animeSchedule := AnimeSchedule{Days: animeCollection}
	for _, anime := range animes {
		for idx, collection := range animeSchedule.Days {
			if anime.Date == collection.Date {
				animeSchedule.Days[idx].Animes = append(animeSchedule.Days[idx].Animes, anime)
			}
		}
	}

	return animeSchedule, nil
}

func addZeroTo2DigitNum(num int) string {
	if num > 9 {
		return fmt.Sprintf("%d", num)
	}
	return fmt.Sprintf("0%d", num)
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
