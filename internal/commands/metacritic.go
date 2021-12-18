package commands

import (
	"fmt"
	"github.com/Mispon/stewart-bot/internal/config"
	"github.com/Mispon/stewart-bot/internal/utils"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type MetacriticProcessor struct{}

// Check checks if a module needs to be executed
func (p MetacriticProcessor) Check(message *discordgo.MessageCreate, _ bool) bool {
	return utils.HasAnyOf(message.Content, []string{"что нового"})
}

// Execute runs module logic
func (p MetacriticProcessor) Execute(message *discordgo.MessageCreate, session *discordgo.Session) {
	cfg := config.GetConfig()

	rssUrl := getRSSUrl(message.Content, cfg)
	if len(rssUrl) == 0 {
		_, _ = session.ChannelMessageSend(message.ChannelID, "Уточни где? В играх, в фильмах?")
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		log.Error().Err(err).Str("metacritic", "failed to parse RSS url").Send()
		return
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s\n\n", strings.TrimSpace(feed.Title)))
	for i, item := range feed.Items[:5] {
		pubTime, _ := time.Parse(time.RFC1123Z, item.Published)

		sb.WriteString(fmt.Sprintf("%v. **%s**\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("%v\n", cutDescription(item.Description)))
		sb.WriteString(fmt.Sprintf("%v\n", item.Link))
		sb.WriteString(fmt.Sprintf("%v\n\n", pubTime.Format("2006-01-02")))
	}

	_, err = session.ChannelMessageSend(message.ChannelID, sb.String())
	if err != nil {
		log.Error().Err(err).Str("metacritic", "failed to send message to channel").Send()
	}
}

func cutDescription(itemDesc string) string {
	desc := strings.TrimSpace(itemDesc)
	if len(desc) > 250 {
		return desc[:250] + "..."
	}
	return desc
}

func getRSSUrl(message string, cfg *config.Config) string {
	if strings.Contains(message, "игр") {
		return cfg.Metacritic.GamesUrl
	}

	if utils.HasAnyOf(message, []string{"кино", "фильм"}) {
		return cfg.Metacritic.MoviesUrl
	}

	return ""
}
