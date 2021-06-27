package commands

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"stewart-bot/app/config"
	"stewart-bot/app/utils"
)

type JokeProcessor struct {}

// Check checks if a module needs to be executed
func (p JokeProcessor) Check(message *discordgo.MessageCreate, wasAsked bool) bool {
	cfg := config.GetConfig()
	return wasAsked && utils.HasAnyOf(message.Content, cfg.Commands.Joke)
}

// Execute runs module logic
func (p JokeProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session)  {
	cfg := config.GetConfig()

	res, err := utils.MakeHTTPRequest(cfg.JokeUrl)
	if err != nil {
		log.Printf("[WARN] failed to make request %s, error=%v", cfg.JokeUrl, err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[WARN] failed to read body, error=%v", err)
		return
	}

	text := string(body)
	text = strings.TrimPrefix(text, `{"content":"`)
	text = strings.TrimSuffix(text, `"}`)

	tr := transform.NewReader(strings.NewReader(text), charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		log.Printf("[WARN] failed to convert string to utf, error=%v", err)
		return
	}

	session.ChannelMessageSend(message.ChannelID, string(buf))
}
