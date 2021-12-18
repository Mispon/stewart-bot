package commands

import (
	"github.com/Mispon/stewart-bot/internal/config"
	"github.com/Mispon/stewart-bot/internal/utils"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageProcessor interface {
	Check(message *discordgo.MessageCreate, askedMe bool) bool
	Execute(message *discordgo.MessageCreate, session *discordgo.Session)
}

var processors = []MessageProcessor{
	PingProcessor{},
	ClearProcessor{},
	JokeProcessor{},
	QuoteProcessor{},
	HoroscopeProcessor{},
	MetacriticProcessor{},
}

// OnMessage Handle chat commands
func OnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	message.Content = prepareContent(message.Content)
	askedMe := askedMe(message.Content)

	for _, processor := range processors {
		if processor.Check(message, askedMe) {
			processor.Execute(message, session)
		}
	}
}

// prepareContent clears message from unimportant symbols
func prepareContent(content string) string {
	reg := regexp.MustCompile(`[.,/#!$%^&*;?:{}=-_~()]`)
	return reg.ReplaceAllString(content, "")
}

// askedMe returns true if message contains any of bot names
func askedMe(content string) bool {
	cfg := config.GetConfig()
	words := strings.Split(content, " ")

	for _, word := range words {
		if utils.IndexOf(word, cfg.BotNames) >= 0 {
			return true
		}
	}
	return false
}
