package config

import (
	"github.com/mispon/stewart-bot/internal/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Config struct {
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
	}
}

// ReadConfig reads bot settings from yaml config
func ReadConfig(name string) (config *Config, err error) {
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

	return config, nil
}
