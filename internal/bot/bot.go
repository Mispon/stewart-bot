package bot

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/commands"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type Bot struct {
	config   *config.Config
	token    string
	commands []commands.Command
	session  *discordgo.Session
}

// New creates new bot instance
func New(cfg *config.Config, token string) *Bot {
	bot := &Bot{
		config: cfg,
		token:  token,
		commands: []commands.Command{
			&commands.PingCommand{},
			&commands.ClearCommand{},
			&commands.JokeCommand{},
			&commands.QuoteCommand{},
			&commands.HoroscopeCommand{},
			&commands.MetacriticCommand{},
		},
	}

	for _, cmd := range bot.commands {
		cmd.WithConfig(cfg)
	}

	return bot
}

// Run bot
func (b *Bot) Run() error {
	discord, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return err
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.AddHandler(b.onMessage)

	err = discord.Open()
	if err != nil {
		return err
	}

	b.session = discord

	logrus.Infof("Stewart-bot v%s successfully started!", b.config.Version)
	return nil
}

func (b *Bot) Close() {
	utils.Close(b.session.Close)
}

// onMessage Handle chat commands
func (b *Bot) onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	message.Content = prepareContent(message.Content)
	askedMe := b.askedMe(message.Content)

	for _, cmd := range b.commands {
		if cmd.Check(message, askedMe) {
			cmd.Execute(message, session)
		}
	}
}

// prepareContent clears message from unimportant symbols
func prepareContent(content string) string {
	reg := regexp.MustCompile(`[.,/#!$%^&*;?:{}=-_~()]`)
	return reg.ReplaceAllString(content, "")
}

// askedMe returns true if message contains any of bot names
func (b *Bot) askedMe(content string) bool {
	words := strings.Split(content, " ")

	for _, word := range words {
		if utils.IndexOf(word, b.config.BotNames) >= 0 {
			return true
		}
	}
	return false
}
