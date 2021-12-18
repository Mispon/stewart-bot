package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type ClearCommand struct{}

// Check checks if a module needs to be executed
func (p ClearCommand) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, []string{"clear", "клир"})
}

// Execute runs module logic
func (p ClearCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	messages, err := session.ChannelMessages(message.ChannelID, 0, "", "", "")
	if err != nil {
		logrus.Error("failed to get messages from channel")
	}

	var messagesIds []string
	for _, msg := range messages {
		messagesIds = append(messagesIds, msg.ID)
	}

	err = session.ChannelMessagesBulkDelete(message.ChannelID, messagesIds)
	if err != nil {
		logrus.Error("failed to bulk delete all messages")
	}
}

// WithConfig setup config pointer
func (p *ClearCommand) WithConfig(*config.Config) {
	// don't use
}
