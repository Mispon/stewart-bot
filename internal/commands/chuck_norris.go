package commands

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/utils"
	"github.com/sirupsen/logrus"
)

type (
	chuckCommand struct {
		config *config.Config
	}

	jokeBody struct {
		Value string `json:"value"`
	}
)

// NewChuckCmd creates new instance
func NewChuckCmd(config *config.Config) Command {
	return &chuckCommand{
		config: config,
	}
}

// Check checks if a module needs to be executed
func (p chuckCommand) Check(message *discordgo.MessageCreate, askedMe bool) bool {
	return askedMe && utils.HasAnyOf(message.Content, p.config.Commands.Chuck.Triggers)
}

// Execute runs module logic
func (p chuckCommand) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	res, err := utils.SendGet(p.config.ChuckNorrisUrl)
	if err != nil {
		logrus.WithField("command", "chuck").Errorf("failed to make request %s", p.config.ChuckNorrisUrl)
		return
	}
	defer utils.Close(res.Body.Close)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithField("command", "chuck").Error("failed to read response body")
		return
	}

	var joke jokeBody
	if err = json.Unmarshal(body, &joke); err != nil {
		logrus.WithField("command", "chuck").Error("failed to parse response body")
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, joke.Value)
	if err != nil {
		logrus.WithField("command", "chuck").Error("failed to send message to channel")
	}
}
