package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type pingCommand struct {
	config *config.Config
}

// NewPingCmd creates new instance
func NewPingCmd(config *config.Config) Command {
	return &pingCommand{
		config: config,
	}
}

// Check checks if a module needs to be executed
func (p pingCommand) Check(message *discordgo.MessageCreate, _ bool) bool {
	return utils.HasAnyOf(message.Content, p.config.Commands.Ping.Triggers)
}

// Execute runs module logic
func (p pingCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	var answer string
	if strings.Contains(message.Content, "доложить") {
		answer = "Служу KR4K3Nу!"
	} else {
		answer = "Так точно!"
	}
	_, _ = session.ChannelMessageSend(message.ChannelID, answer)
}
