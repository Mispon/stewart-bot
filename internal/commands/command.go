package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Check(message *discordgo.MessageCreate, askedMe bool) bool
	Execute(message *discordgo.MessageCreate, session *discordgo.Session)
}
