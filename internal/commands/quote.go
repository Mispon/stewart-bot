package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type quoteCommand struct {
	config *config.Config
}

// NewQuoteCmd creates new instance
func NewQuoteCmd(config *config.Config) Command {
	return &quoteCommand{
		config: config,
	}
}

// Check checks if a module needs to be executed
func (p quoteCommand) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	return wasAsked && utils.HasAnyOf(message.Content, p.config.Commands.Quote.Triggers)
}

// Execute runs module logic
func (p quoteCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	res, err := utils.SendGet(p.config.QuoteUrl)
	if err != nil {
		logrus.WithField("command", "quote").Errorf("failed to make request %s", p.config.QuoteUrl)
		return
	}
	defer utils.Close(res.Body.Close)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithField("command", "quote").Error("failed to read response body")
		return
	}

	quote := struct {
		QuoteText   string
		QuoteAuthor string
	}{}
	if err = json.Unmarshal(body, &quote); err != nil {
		logrus.WithField("command", "quote").Error("failed to deserialize json body")
	}

	text := fmt.Sprintf(`%s`, quote.QuoteText)
	if len(quote.QuoteAuthor) > 0 {
		text += fmt.Sprintf("\n%s", quote.QuoteAuthor)
	}

	_, err = session.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		logrus.WithField("command", "quote").Error("failed to send message to channel")
	}
}
