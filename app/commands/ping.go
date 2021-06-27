package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Mispon/stewart-bot/app/config"
	"github.com/Mispon/stewart-bot/app/utils"
)

type PingProcessor struct{}

// Check checks if a module needs to be executed
func (p PingProcessor) Check(message *discordgo.MessageCreate, _ bool) bool {
	return utils.HasAnyOf(message.Content, []string{"ping", "pong"})
}

// Execute runs module logic
func (p PingProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	cfg := config.GetConfig()
	var answer string
	if message.Content == "ping" {
		answer = fmt.Sprintf("%s pong!", cfg.Author)
	} else {
		answer = fmt.Sprintf("%s ping!", cfg.Author)
	}
	_, _ = session.ChannelMessageSend(message.ChannelID, answer)
}
