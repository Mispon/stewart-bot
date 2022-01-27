package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/mispon/stewart-bot/internal/utils"
)

type (
	BotCommand struct {
		Triggers []string `yaml:"triggers"`
		Info     string   `yaml:"info"`
	}

	Config struct {
		Author   string   `yaml:"author"`
		Version  string   `yaml:"version"`
		BotNames []string `yaml:"bot_names"`
		Commands struct {
			Joke       BotCommand `yaml:"joke"`
			Quote      BotCommand `yaml:"quote"`
			DudeQuote  BotCommand `yaml:"dude_quote"`
			Horoscope  BotCommand `yaml:"horoscope"`
			Chuck      BotCommand `yaml:"chuck"`
			Metacritic BotCommand `yaml:"metacritic"`
			Thanks     BotCommand `yaml:"thanks"`
			Ping       BotCommand `yaml:"ping"`
			Help       BotCommand `yaml:"help"`
		} `yaml:"commands"`
		JokeUrl        string `yaml:"joke_url"`
		QuoteUrl       string `yaml:"quote_url"`
		HoroscopeUrl   string `yaml:"horoscope_url"`
		ChuckNorrisUrl string `yaml:"chuck_norris_url"`
		BalabobaUrl    string `yaml:"balaboba_url"`
		Metacritic     struct {
			GamesUrl  string `yaml:"games_url"`
			MoviesUrl string `yaml:"movies_url"`
		} `yaml:"metacritic"`
		Members []struct {
			ID     string `yaml:"id"`
			Name   string `yaml:"name"`
			Zodiac string `yaml:"zodiac"`
		} `yaml:"members"`
		Options Options
	}

	Options struct {
		ServerID       string
		MainChannelID  string
		VoiceChannelID string
	}

	OptionsFn func(o *Options)
)

// ReadConfig reads bot settings from yaml config
func ReadConfig(name string, opts ...OptionsFn) (config *Config, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer utils.Close(file.Close)

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	for _, optFn := range opts {
		optFn(&config.Options)
	}

	return config, nil
}

func WithServerID(serverID string) OptionsFn {
	return func(o *Options) {
		o.ServerID = serverID
	}
}

func WithMainChannelID(mainChannelID string) OptionsFn {
	return func(o *Options) {
		o.MainChannelID = mainChannelID
	}
}

func WithVoiceChannelID(voiceChannelID string) OptionsFn {
	return func(o *Options) {
		o.VoiceChannelID = voiceChannelID
	}
}

func (b BotCommand) String() string {
	return fmt.Sprintf("\t**[%s]** - %s\n", strings.Join(b.Triggers, ", "), b.Info)
}
