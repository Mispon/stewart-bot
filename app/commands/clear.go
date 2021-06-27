package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/Mispon/stewart-bot/app/utils"
)

type ClearProcessor struct{}

// Check checks if a module needs to be executed
func (p ClearProcessor) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, []string{"clear", "клир"})
}

// Execute runs module logic
func (p ClearProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	messages, err := session.ChannelMessages(message.ChannelID, 0, "", "", "")
	if err != nil {
		log.Error().Err(err).Msg("failed to get messages from channel")
	}

	var messagesIds []string
	for _, msg := range messages {
		messagesIds = append(messagesIds, msg.ID)
	}

	err = session.ChannelMessagesBulkDelete(message.ChannelID, messagesIds)
	if err != nil {
		log.Error().Err(err).Msg("failed to bulk delete all messages")
	}
}
