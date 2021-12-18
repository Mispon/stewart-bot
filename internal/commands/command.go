package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mispon/stewart-bot/internal/config"
)

type Command interface {
	Check(message *discordgo.MessageCreate, askedMe bool) bool
	Execute(message *discordgo.MessageCreate, session *discordgo.Session)
	WithConfig(cfg *config.Config)
}
