package static

import "fmt"

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}

// Commands a map that contains all commands with a description
var Commands = map[string]Command{
	"help": {
		Description:        "Shows information about all commands",
		SpecialPermissions: "None",
	},
	"ping": {
		Description:        "Send 'Pong!' back to the user",
		SpecialPermissions: "None",
	},
	"animeSites": {
		Description:        "Links the site https://piracy.moe as an anime streaming site index",
		SpecialPermissions: "None",
	},
	"bestReleases": {
		Description:        "Links the site https://releases.moe as an overview site for best anime torrents",
		SpecialPermissions: "None",
	},
	"nhentai": {
		Description:        fmt.Sprintf("%s123456 or %s12345 returns the complete link to the nhentai page", Prefix, Prefix),
		SpecialPermissions: "None",
	},
}

// AnimeNewsChannelID id
const AnimeNewsChannelID = "805446666428612659"

// HentaiNewsChannelID id
const HentaiNewsChannelID = "809693883809923082"

// Prefix for discord commands
const Prefix = "+"

// AnimeOfCurrentDay y or n
const AnimeOfCurrentDay = true

const (
	// AmiAmi keyword for amiami.com
	AmiAmi = "amiami"
	// CuddlyOctopus keyword for cuddlyoctopus.com
	CuddlyOctopus = "cuddlyoctopus"
	// NHentai keyword for nhentai.net
	NHentai = "nhentai"
	// SolarisJapan keyword for solarisjapan.com
	SolarisJapan = "SolarisJapan"
	// TwitterProxy keyword for twitter shortend URLs
	TwitterProxy = "TwitterProxy"
)
