package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bwmarrin/discordgo"

	"stewart-bot/app/config"
	"stewart-bot/app/utils"
)

type QuoteProcessor struct {}

// Check checks if a module needs to be executed
func (p QuoteProcessor) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	cfg := config.GetConfig()
	return wasAsked && utils.HasAnyOf(message.Content, cfg.Commands.Quote)
}

// Execute runs module logic
func (p QuoteProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session)  {
	cfg := config.GetConfig()

	res, err := utils.MakeHTTPRequest(cfg.QuoteUrl)
	if err != nil {
		log.Printf("[WARN] failed to make request %s, error=%v", cfg.QuoteUrl, err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[WARN] failed to read body, error=%v", err)
		return
	}

	quote := struct {
		QuoteText 	string
		QuoteAuthor string
	}{}
	json.Unmarshal(body, &quote)

	text := fmt.Sprintf(`*%s*`, quote.QuoteText)
	if len(quote.QuoteAuthor) > 0 {
		text += fmt.Sprintf("\n%s", quote.QuoteAuthor)
	}

	session.ChannelMessageSend(message.ChannelID, text)
}