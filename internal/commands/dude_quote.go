package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/balaboba"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type dudeQuoteCommand struct {
	config *config.Config
	bb     balaboba.Balaboba
}

// NewDudeQuoteCmd creates new instance
func NewDudeQuoteCmd(config *config.Config) Command {
	return &dudeQuoteCommand{
		config: config,
		bb:     balaboba.New(config.BalabobaUrl),
	}
}

// Check checks if a module needs to be executed
func (p dudeQuoteCommand) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	return wasAsked && utils.HasAnyOf(message.Content, p.config.Commands.DudeQuote.Triggers)
}

// Execute runs module logic
func (p dudeQuoteCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	var (
		query = p.parseTheme(message.Content)
		intro = 4
	)

	text, err := p.bb.GetText(query, intro)
	if err != nil {
		logrus.WithField("command", "dude_quote").Error("failed to get text from balaboba")
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		logrus.WithField("command", "dude_quote").Error("failed to send message to channel")
	}
}

// parseTheme gets theme from message
func (p dudeQuoteCommand) parseTheme(message string) string {
	sep := p.config.Commands.DudeQuote.Triggers[0]
	values := strings.Split(message, sep)

	if len(values) < 2 {
		return ""
	}

	return strings.TrimSpace(values[1])
}
