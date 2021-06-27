package config

import (
	"os"

	"gopkg.in/yaml.v2"
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

var config *Config

// ReadConfig reads bot settings from yaml config
func ReadConfig(path string) {
	config = &Config{}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}
}

// GetConfig returns config instance
func GetConfig() *Config {
	return config
}
