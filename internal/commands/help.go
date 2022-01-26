package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
	"github.com/sirupsen/logrus"
)

type helpCommand struct {
	cfg *config.Config
}

// NewHelpCmd creates new instance
func NewHelpCmd(cfg *config.Config) Command {
	return &helpCommand{
		cfg: cfg,
	}
}

// Check checks if a module needs to be executed
func (p helpCommand) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, p.cfg.Commands.Help.Triggers)
}

// Execute runs module logic
func (p helpCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	_, err := session.ChannelMessageSend(message.ChannelID, p.GetHelpMessage())
	if err != nil {
		logrus.Error("failed to bulk delete all messages")
	}
}

func (p helpCommand) GetHelpMessage() string {
	var result string

	result += "Команды:\n"
	result += fmt.Sprintf("имена: %s\n", strings.Join(p.cfg.BotNames, ", "))
	result += p.cfg.Commands.Help.String()
	result += p.cfg.Commands.Ping.String()
	result += p.cfg.Commands.Joke.String()
	result += p.cfg.Commands.Horoscope.String()
	result += p.cfg.Commands.Quote.String()
	result += p.cfg.Commands.DudeQuote.String()
	result += p.cfg.Commands.Metacritic.String()
	result += p.cfg.Commands.Chuck.String()
	result += "имя + сообщение - просто поболтать"

	return result
}
