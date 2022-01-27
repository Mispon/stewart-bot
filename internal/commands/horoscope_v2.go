package commands

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/balaboba"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type horoscopeCommandV2 struct {
	config *config.Config
	bb     balaboba.Balaboba
}

// NewHoroscopeV2Cmd creates new instance of command
func NewHoroscopeV2Cmd(config *config.Config) JobCommand {
	return &horoscopeCommandV2{
		config: config,
		bb:     balaboba.New(config.BalabobaUrl),
	}
}

// TriggerTime returns job execution time
func (p horoscopeCommandV2) TriggerTime() (h int, m int) {
	return 6, 30 // UTC
}

// Run executes command from job
func (p horoscopeCommandV2) Run(session *discordgo.Session) error {
	ri := rand.Intn(len(p.config.Members))
	member := p.config.Members[ri]

	horoscope, ok := p.getPersonalHoroscope(member.Zodiac)
	if !ok {
		return errors.New(fmt.Sprintf("received empty horoscope for user %s", member.ID))
	}

	message := fmt.Sprintf("<@%s>, %s\n%s", member.ID, member.Zodiac, horoscope)
	_, err := session.ChannelMessageSend(p.config.Options.MainChannelID, message)

	return err
}

// Check checks if a module needs to be executed
func (p horoscopeCommandV2) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, p.config.Commands.Horoscope.Triggers)
}

// Execute runs module logic
func (p horoscopeCommandV2) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	zodiac, ok := p.getUserZodiac(message.Author.ID)
	if !ok {
		_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("user %s not found in config members list", message.Author.ID))
		return
	}

	horoscope, ok := p.getPersonalHoroscope(zodiac)
	if !ok {
		_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(`horoscope for "%s" not found`, zodiac))
		return
	}

	response := fmt.Sprintf("<@%s>, %s\n%s", message.Author.ID, zodiac, horoscope)
	_, err := session.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		logrus.WithField("command", "horoscope").Error("horoscope", "failed to send message to channel")
	}
}

func (p horoscopeCommandV2) getPersonalHoroscope(zodiac string) (string, bool) {
	var (
		query = zodiac
		intro = 10
	)

	text, err := p.bb.GetText(query, intro)
	if err != nil {
		logrus.WithField("command", "horoscope").Error("failed to get text from balaboba")
		return "", false
	}

	return text, true
}

func (p *horoscopeCommandV2) getUserZodiac(authorID string) (string, bool) {
	for _, member := range p.config.Members {
		if member.ID == authorID {
			return member.Zodiac, true
		}
	}

	return "", false
}
