package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/request"
	"github.com/stegosawr/Albedo/static"
)

// AnimeCard = 1 Anime
type AnimeCard struct {
	URL      string
	Title    string
	ThumbURL string
	Episode  string
	Date     string
}

// AnimeCollection Grouped by day
type AnimeCollection struct {
	Date   string
	Animes []AnimeCard
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
				s.ChannelMessageSend(channel, fmt.Sprintf("No releases today! Check release schedule here %s", fmt.Sprintf("https://www.underhentai.net/releases-%d/", time.Now().Year())))
			}
		}
		for _, day := range schedule.Days {
			for _, a := range day.Animes {
				time.Sleep(2 * time.Second)
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

func scraperAnime() (AnimeSchedule, error) {
	header := map[string]string{
		"cookie": `preferences=%7B%22time_zone%22%3A%22Europe%2FBerlin%22%2C%22sortby%22%3A%22popularity%22%2C%22titles%22%3A%22romaji%22%2C%22ongoing%22%3A%22all%22%2C%22use_24h_clock%22%3Afalse%2C%22night_mode%22%3Atrue%2C%22reveal_spoilers%22%3Atrue%7D;`,
	}

	resp, err := request.Request(http.MethodGet, "https://www.livechart.me/timetable", header)
	if err != nil {
		return AnimeSchedule{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AnimeSchedule{}, err
	}
	htmlString := string(body)
	today := time.Now()
	startDate := time.Date(today.Year(), today.Month(), today.Day(), -1, 0, 0, 0, time.UTC).Unix()
	endDate := time.Date(today.Year(), today.Month(), today.Day()+1, -1, 0, 0, 0, time.UTC).Unix()
	re := regexp.MustCompile(fmt.Sprintf(`data-timetable-day-start="%d[\s\S]*?%d`, startDate, endDate))
	matchedDay := re.FindString(htmlString)

	re = regexp.MustCompile(`class="lazy-img".+?alt="([^"]*)".*?"(.+?(\d[^/]*)[^"]*)`)
	matchedAnimeInfo := re.FindAllStringSubmatch(matchedDay, -1) // 1=Title 2=PosterURL 3=AnimeID

	re = regexp.MustCompile(`</span>([^<]+)</div>`)
	matchedEPDescriptions := re.FindAllStringSubmatch(matchedDay, -1) // 1=Episode info

	if len(matchedAnimeInfo) != len(matchedEPDescriptions) {
		return AnimeSchedule{}, fmt.Errorf("Internal Error: found %d anime and %d EP descriptions", len(matchedAnimeInfo), len(matchedEPDescriptions))
	}

	animes := []AnimeCard{}
	for idx, matchedImgTag := range matchedAnimeInfo {
		animes = append(animes, AnimeCard{
			URL:      fmt.Sprintf("%s%s", "https://www.livechart.me/anime/", matchedImgTag[3]),
			Title:    matchedImgTag[1],
			ThumbURL: matchedImgTag[2],
			Episode:  matchedEPDescriptions[idx][1],
		})
	}

	return AnimeSchedule{[]AnimeCollection{
		{
			Date:   today.String(),
			Animes: animes,
		},
	}}, nil
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
