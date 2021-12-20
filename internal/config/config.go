package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"

	"github.com/mispon/stewart-bot/internal/utils"
)

type (
	Config struct {
		Author   string   `yaml:"author"`
		Version  string   `yaml:"version"`
		BotNames []string `yaml:"bot_names"`
		Commands struct {
			Joke      []string `yaml:"joke"`
			Quote     []string `yaml:"quote"`
			Horoscope []string `yaml:"horoscope"`
		} `yaml:"commands"`
		JokeUrl      string `yaml:"joke_url"`
		QuoteUrl     string `yaml:"quote_url"`
		HoroscopeUrl string `yaml:"horoscope_url"`
		Metacritic   struct {
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
	path, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
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
