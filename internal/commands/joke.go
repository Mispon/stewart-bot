package commands

import (
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
)

type jokeCommand struct {
	config *config.Config
}

// NewJokeCmd creates new instance
func NewJokeCmd(config *config.Config) Command {
	return &jokeCommand{
		config: config,
	}
}

// Check checks if a module needs to be executed
func (p jokeCommand) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	return wasAsked && utils.HasAnyOf(message.Content, p.config.Commands.Joke.Triggers)
}

// Execute runs module logic
func (p jokeCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	res, err := utils.SendGet(p.config.JokeUrl)
	if err != nil {
		logrus.WithField("command", "joke").Errorf("failed to make request %s", p.config.JokeUrl)
		return
	}
	defer utils.Close(res.Body.Close)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithField("command", "joke").Error("failed to read response body")
		return
	}

	text := string(body)
	text = strings.TrimPrefix(text, `{"content":"`)
	text = strings.TrimSuffix(text, `"}`)

	tr := transform.NewReader(strings.NewReader(text), charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		logrus.WithField("command", "joke").Error("failed to convert string to utf")
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, string(buf))
	if err != nil {
		logrus.WithField("command", "joke").Error("failed to send message to channel")
	}
}
