package commands

import (
	"errors"
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

// NewHoroscopeCmd creates new instance of command
func NewHoroscopeCmd(config *config.Config) JobCommand {
	return &horoscopeCommand{
		config: config,
	}
}

// TriggerTime returns job execution time
func (p horoscopeCommand) TriggerTime() (h int, m int) {
	return 9, 30
}

// Run executes command from job
func (p horoscopeCommand) Run(session *discordgo.Session) error {
	ri := rand.Intn(len(p.config.Members))
	member := p.config.Members[ri]

	horoscope, ok := p.getPersonalHoroscope(p.config.HoroscopeUrl, member.Zodiac)
	if !ok {
		return errors.New(fmt.Sprintf("received empty horoscope for user %s", member.ID))
	}

	message := fmt.Sprintf("<@%s>,%s", member.ID, horoscope)
	_, err := session.ChannelMessageSend(p.config.Options.MainChannelID, message)

	return err
}

// Check checks if a module needs to be executed
func (p horoscopeCommand) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, p.config.Commands.Horoscope.Triggers)
}

// Execute runs module logic
func (p horoscopeCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	zodiac, ok := p.getUserZodiac(message.Author.ID)
	if !ok {
		_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("user %s not found in config members list", message.Author.ID))
		return
	}

	horoscope, ok := p.getPersonalHoroscope(p.config.HoroscopeUrl, zodiac)
	if !ok {
		_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(`horoscope for "%s" not found`, zodiac))
		return
	}

	_, err := session.ChannelMessageSend(message.ChannelID, horoscope.String())
	if err != nil {
		logrus.WithField("command", "horoscope").Error("horoscope", "failed to send message to channel")
	}
}

// getPersonalHoroscope returns personal horoscope by zodiac
func (p horoscopeCommand) getPersonalHoroscope(url, zodiac string) (*horoscopeItem, bool) {
	full := p.getFullHoroscope(url)

	var result horoscopeItem
	for _, h := range full {
		if h.Header == zodiac {
			result = h
			break
		}
	}

	if len(result.Content) == 0 {
		logrus.WithField("command", "horoscope").Warnf(`horoscope for "%s" not found!`, zodiac)
		return nil, false
	}

	return &result, true
}

// getFullHoroscope returns list of daily horoscopes
func (p horoscopeCommand) getFullHoroscope(url string) (horoscope []horoscopeItem) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		logrus.WithField("command", "horoscope").Error("failed to get doc")
	}

	entry := htmlquery.FindOne(doc, "//div[@class='entry']")
	headers := htmlquery.Find(entry, "//h4")
	contents := htmlquery.Find(entry, "//p")
	count := len(headers) - 1

	horoscope = make([]horoscopeItem, count)
	for i := 0; i < count; i++ {
		hi := horoscopeItem{
			Header:  strings.TrimSpace(htmlquery.InnerText(headers[i])),
			Content: strings.TrimSpace(htmlquery.InnerText(contents[i])),
		}
		horoscope[i] = hi
	}

	return horoscope
}

func (p *horoscopeCommand) getUserZodiac(authorID string) (string, bool) {
	for _, member := range p.config.Members {
		if member.ID == authorID {
			return member.Zodiac, true
		}
	}

	return "", false
}

type horoscopeItem struct {
	Header  string
	Content string
}

func (h horoscopeItem) String() string {
	return fmt.Sprintf("**%s**\n%s", strings.TrimSpace(h.Header), strings.TrimSpace(h.Content))
}
