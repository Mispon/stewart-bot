package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type thanksCommand struct {
	config *config.Config
}

// NewThanksCmd creates new instance
func NewThanksCmd(config *config.Config) Command {
	return &thanksCommand{
		config: config,
	}
}

// Check checks if a module needs to be executed
func (p thanksCommand) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	return wasAsked && utils.HasAnyOf(message.Content, p.config.Commands.Thanks.Triggers)
}

// Execute runs module logic
func (p thanksCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	_, _ = session.ChannelMessageSend(message.ChannelID, "Служу KR4K3Nу!")
}
