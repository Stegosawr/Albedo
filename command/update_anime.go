package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stegosawr/Albedo/static"
)

// MediaInfo = 1 Anime,Hentai,Manga
type MediaInfo struct {
	URL           string
	Title         string
	Description   string
	Categories    []string
	ThumbURL      string
	Episode       string
	AiringAt      string
	Type          string
	AuthorIconURL string
}

type tag struct {
	Name string `json:"name,omitempty"`
}

type title struct {
	Romaji  string `json:"romaji,omitempty"`
	English string `json:"english,omitempty"`
	Native  string `json:"native,omitempty"`
}

type media struct {
	ID          uint32   `json:"id,omitempty"`
	IsAdult     bool     `json:"isAdult,omitempty"`
	Title       title    `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	SiteURL     string   `json:"siteUrl,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Tags        []tag    `json:"tags,omitempty"`
}

type airingSchedule struct {
	Media    media  `json:"media,omitempty"`
	Episode  uint32 `json:"episode,omitempty"`
	AiringAt int64  `json:"airingAt,omitempty"`
}

type page struct {
	AiringSchedules []airingSchedule `json:"airingSchedules,omitempty"`
}

type data struct {
	Page page `json:"page"`
}

type animeAiring struct {
	Data data `json:"data"`
}

type graphqlReq struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

const aniListApi = "https://graphql.anilist.co"
const aniListMediaThumb = "https://img.anili.st/media/"

// UpdateAnime Anime Schedule
func UpdateAnime(s *discordgo.Session) {

	_ = DeleteAllMessagesInChannel(s, static.AnimeNewsChannelID)

	_ = DeleteAllMessagesInChannel(s, static.HentaiNewsChannelID)

	mediaOfToday, err := scraperAnime()
	if err != nil {
		s.ChannelMessageSend(static.AnimeNewsChannelID, err.Error())
		return
	}

	//s.ChannelMessageSend(static.AnimeNewsChannelID, "--------------------------\n  Anime of Today\n--------------------------")
	//s.ChannelMessageSend(static.HentaiNewsChannelID, "--------------------------\n  Hentai of Today\n--------------------------")
	hasAdult := false
	for _, v := range mediaOfToday {
		if v.Type == "Hentai" {
			hasAdult = true
			break
		}
	}

	if !hasAdult {
		s.ChannelMessageSend(static.HentaiNewsChannelID, fmt.Sprintf("No releases today! Check release schedule here %s", fmt.Sprintf("https://www.underhentai.net/releases-%d/", time.Now().Year())))
	}

	channel := static.AnimeNewsChannelID
	for _, m := range mediaOfToday {
		switch m.Type {
		case "Anime":
			channel = static.AnimeNewsChannelID
		case "Hentai":
			channel = static.HentaiNewsChannelID
		}
		time.Sleep(5 * time.Second)
		s.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
			URL:         m.URL,
			Title:       m.Title,
			Description: appendStringSeq("\n", m.Description, appendStringSeq(",", m.Categories...)),
			Image: &discordgo.MessageEmbedImage{
				URL: m.ThumbURL,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    appendStringSeq(" - ", m.Type, m.Episode, m.AiringAt),
				IconURL: m.AuthorIconURL,
			},
		})
		//s.ChannelMessageSend(channel, m.ThumbURL)
	}

}

func scraperAnime() ([]MediaInfo, error) {
	//today := time.Date(2021, 5, 27, 1, 10, 10, 0, time.FixedZone("CEST", 2*60*60))
	today := time.Now()
	startDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).Unix()
	endDate := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location()).Unix()

	query := `
	query($airingAt_greater: Int, $airingAt_lesser: Int){
		Page(page: 1) {
		  airingSchedules(airingAt_greater: $airingAt_greater, airingAt_lesser: $airingAt_lesser, sort: TIME) {
			media {
			  id
			  isAdult
			  title {
				romaji
				english
				native
			  }
			  description
			  siteUrl
			  coverImage{
				color
			  }
			  genres
			  tags{name}
			}
			airingAt
			episode
		  }
		}
	}	  
	`

	variables := fmt.Sprintf("{ \"airingAt_greater\": %d, \"airingAt_lesser\": %d }", startDate, endDate)

	graphqlBody := &graphqlReq{Query: query, Variables: variables}
	reqBody, err := json.Marshal(graphqlBody)
	if err != nil {
		return nil, err
	}

	jsonRes, err := http.Post(aniListApi, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer jsonRes.Body.Close()

	resBody, err := ioutil.ReadAll(jsonRes.Body)
	if err != nil {
		return nil, err
	}

	airingAnime := &animeAiring{}
	err = json.Unmarshal(resBody, &airingAnime)
	if err != nil {
		return nil, err
	}

	animes := []MediaInfo{}
	for _, schedule := range airingAnime.Data.Page.AiringSchedules {
		mediaType := "Anime"
		if schedule.Media.IsAdult {
			mediaType = "Hentai"
		}

		t := time.Unix(schedule.AiringAt, 0)

		title := schedule.Media.Title.Romaji
		if schedule.Media.Title.English != "" {
			title = title + " / " + schedule.Media.Title.English
		}
		if schedule.Media.Title.Native != "" {
			title = title + " / " + schedule.Media.Title.Native
		}

		zone, offset := t.Zone()

		//fix description
		re := regexp.MustCompile(`</?[ib]>`)
		schedule.Media.Description = re.ReplaceAllString(schedule.Media.Description, "")
		firstHtmlTag := strings.Index(schedule.Media.Description, "<")
		if firstHtmlTag == -1 {
			firstHtmlTag = len(schedule.Media.Description)
		}
		schedule.Media.Description = schedule.Media.Description[:firstHtmlTag] + fmt.Sprintf("[(read more)](%s)", schedule.Media.SiteURL)

		categories := []string{}
		switch mediaType {
		case "Anime":
			categories = wrappStringIn("***", schedule.Media.Genres...)
		case "Hentai":
		}
		if mediaType == "Hentai" {
			for _, tag := range schedule.Media.Tags {
				categories = append(categories, wrappStringIn("***", tag.Name)...)
			}
		}

		animes = append(animes, MediaInfo{
			URL:           schedule.Media.SiteURL,
			Title:         title,
			Description:   schedule.Media.Description,
			Categories:    categories,
			Episode:       "Episode " + fmt.Sprint(schedule.Episode),
			ThumbURL:      fmt.Sprintf("%s%d", aniListMediaThumb, schedule.Media.ID),
			AiringAt:      fmt.Sprintf("Airing at %s %s+%d", t.Local().Format("15:04:05"), zone, offset/60/60),
			Type:          mediaType,
			AuthorIconURL: "https://anilist.co/img/icons/favicon-32x32.png",
		})
	}

	return animes, nil
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

func appendStringSeq(sep string, pieces ...string) string {
	out := ""
	for _, piece := range pieces {
		if out == "" {
			out = piece
			continue
		}
		out = out + sep + piece
	}
	return out
}

func wrappStringIn(c string, elements ...string) []string {
	out := []string{}
	for _, elem := range elements {
		out = append(out, c+elem+c)
	}
	return out
}
