package commands

import (
	"github.com/bwmarrin/discordgo"
	"stewart-bot/app/utils"
)

type ClearProcessor struct {}

// Check checks if a module needs to be executed
func (p ClearProcessor) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, []string {"clear", "клир"})
}

// Execute runs module logic
func (p ClearProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session)  {
	messages, _ := session.ChannelMessages(message.ChannelID, 0, "", "", "")

	var messagesIds []string
	for _, msg := range messages {
		messagesIds = append(messagesIds, msg.ID)
	}

	session.ChannelMessagesBulkDelete(message.ChannelID, messagesIds)
}