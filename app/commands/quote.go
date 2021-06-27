package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/Mispon/stewart-bot/app/config"
	"github.com/Mispon/stewart-bot/app/utils"
)

type QuoteProcessor struct{}

// Check checks if a module needs to be executed
func (p QuoteProcessor) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	cfg := config.GetConfig()
	return wasAsked && utils.HasAnyOf(message.Content, cfg.Commands.Quote)
}

// Execute runs module logic
func (p QuoteProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	cfg := config.GetConfig()

	res, err := utils.MakeHTTPRequest(cfg.QuoteUrl)
	if err != nil {
		log.Error().Err(err).Msgf("failed to make request %s", cfg.QuoteUrl)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("quote", "failed to read response body").Send()
		return
	}

	quote := struct {
		QuoteText   string
		QuoteAuthor string
	}{}
	if err = json.Unmarshal(body, &quote); err != nil {
		log.Error().Err(err).Str("quote", "failed to deserialize json body").Send()
	}

	text := fmt.Sprintf(`*%s*`, quote.QuoteText)
	if len(quote.QuoteAuthor) > 0 {
		text += fmt.Sprintf("\n%s", quote.QuoteAuthor)
	}

	_, err = session.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		log.Error().Err(err).Str("quote", "failed to send message to channel").Send()
	}
}
