package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/balaboba"
	"github.com/mispon/stewart-bot/internal/config"
)

type dialogueCommand struct {
	config *config.Config
	bb     balaboba.Balaboba
}

// NewDialogueCmd creates new instance
func NewDialogueCmd(config *config.Config) Command {
	return &dialogueCommand{
		config: config,
		bb:     balaboba.New(config.BalabobaUrl),
	}
}

// Check checks if a module needs to be executed
func (p dialogueCommand) Check(_ *discordgo.MessageCreate, wasAsked bool) bool {
	return wasAsked
}

// Execute runs module logic
func (p dialogueCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	var (
		query = p.parseTheme(message.Content)
		intro = 6
	)

	text, err := p.bb.GetText(query, intro)
	if err != nil {
		logrus.WithField("command", "dialogue").Error("failed to get text from balaboba")
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		logrus.WithField("command", "dialogue").Error("failed to send message to channel")
	}
}

// parseTheme gets theme from message
func (p dialogueCommand) parseTheme(message string) string {
	for _, name := range p.config.BotNames {
		message = strings.TrimPrefix(message, name)
	}
	message = strings.TrimPrefix(message, ",")

	return strings.TrimSpace(message)
}
