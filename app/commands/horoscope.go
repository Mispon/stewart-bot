package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/Mispon/stewart-bot/app/config"
	"github.com/Mispon/stewart-bot/app/utils"
)

type HoroscopeProcessor struct{}

type HoroscopeItem struct {
	Header  string
	Content string
}

func (h HoroscopeItem) String() string {
	return fmt.Sprintf("**%v**\n%v", strings.TrimSpace(h.Header), strings.TrimSpace(h.Content))
}

// Check checks if a module needs to be executed
func (p HoroscopeProcessor) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	cfg := config.GetConfig()
	return askedMe && utils.HasAnyOf(message.Content, cfg.Commands.Horoscope)
}

// Execute runs module logic
func (p HoroscopeProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	cfg := config.GetConfig()

	horoscope := getHoroscope(cfg.HoroscopeUrl)

	ri := rand.Intn(len(horoscope))
	rh := fmt.Sprintf("%v", horoscope[ri])

	_, err := session.ChannelMessageSend(message.ChannelID, rh)
	if err != nil {
		log.Error().Err(err).Str("horoscope", "failed to send message to channel").Send()
	}
}

// getHoroscope returns list of daily horoscopes
func getHoroscope(url string) (horoscope []HoroscopeItem) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error().Err(err).Str("horoscope", "failed to get doc").Send()
	}

	entry := htmlquery.FindOne(doc, "//div[@class='entry']")
	headers := htmlquery.Find(entry, "//h4")
	contents := htmlquery.Find(entry, "//p")
	count := len(headers) - 1

	horoscope = make([]HoroscopeItem, count)
	for i := 0; i < count; i++ {
		hi := HoroscopeItem{
			Header:  htmlquery.InnerText(headers[i]),
			Content: htmlquery.InnerText(contents[i]),
		}
		horoscope[i] = hi
	}

	return horoscope
}
