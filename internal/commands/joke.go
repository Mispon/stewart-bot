package commands

import (
	"github.com/Mispon/stewart-bot/internal/config"
	utils2 "github.com/Mispon/stewart-bot/internal/utils"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/bwmarrin/discordgo"
)

type JokeProcessor struct{}

// Check checks if a module needs to be executed
func (p JokeProcessor) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	cfg := config.GetConfig()
	return wasAsked && utils2.HasAnyOf(message.Content, cfg.Commands.Joke)
}

// Execute runs module logic
func (p JokeProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	cfg := config.GetConfig()

	res, err := utils2.MakeHTTPRequest(cfg.JokeUrl)
	if err != nil {
		log.Error().Err(err).Msgf("failed to make request %s to API", cfg.JokeUrl)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("joke", "failed to read response body").Send()
		return
	}

	text := string(body)
	text = strings.TrimPrefix(text, `{"content":"`)
	text = strings.TrimSuffix(text, `"}`)

	tr := transform.NewReader(strings.NewReader(text), charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		log.Error().Err(err).Str("joke", "failed to convert string to utf").Send()
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, string(buf))
	if err != nil {
		log.Error().Err(err).Str("joke", "failed to send message to channel").Send()
	}
}
