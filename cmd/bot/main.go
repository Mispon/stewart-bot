package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	stewart "github.com/mispon/stewart-bot/internal/bot"
	"github.com/mispon/stewart-bot/internal/config"
)

var (
	debug bool
	token string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "--debug")
	flag.StringVar(&token, "token", "", "--token=my_bot_token")

	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	logrus.Infoln("Initialize bot...")

	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		logrus.Fatal(err)
	}

	bot := stewart.New(cfg, token)
	if err := bot.Run(); err != nil {
		logrus.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, os.Interrupt)
	<-sc

	bot.Close()
}
