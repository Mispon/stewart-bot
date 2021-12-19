package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type horoscopeCommand struct {
	config *config.Config
}

// NewHoroscopeCmd creates new instance
func NewHoroscopeCmd(config *config.Config) Command {
	return &horoscopeCommand{
		config: config,
	}
}

type HoroscopeItem struct {
	Header  string
	Content string
}

func (h HoroscopeItem) String() string {
	return fmt.Sprintf("**%v**\n%v", strings.TrimSpace(h.Header), strings.TrimSpace(h.Content))
}

// Check checks if a module needs to be executed
func (p horoscopeCommand) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, p.config.Commands.Horoscope)
}

// Execute runs module logic
func (p horoscopeCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	horoscope := getHoroscope(p.config.HoroscopeUrl)

	// TODO: returns right horoscope for user
	ri := rand.Intn(len(horoscope))
	rh := fmt.Sprintf("%v", horoscope[ri])

	_, err := session.ChannelMessageSend(message.ChannelID, rh)
	if err != nil {
		logrus.
			WithField("command", "horoscope").
			Error("horoscope", "failed to send message to channel")
	}
}

// getHoroscope returns list of daily horoscopes
func getHoroscope(url string) (horoscope []HoroscopeItem) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		logrus.
			WithField("command", "horoscope").
			Error("failed to get doc")
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
