package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	stewart "github.com/mispon/stewart-bot/internal/bot"
	"github.com/mispon/stewart-bot/internal/config"
	"github.com/mispon/stewart-bot/internal/job"

	"github.com/sirupsen/logrus"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "--debug=true")
	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	logrus.Infoln("Initialize bot...")

	var (
		token          = os.Getenv("STEW_TOKEN")
		serverID       = os.Getenv("STEW_SERVER_ID")
		mainChannelID  = os.Getenv("STEW_MAIN_CH")
		voiceChannelID = os.Getenv("STEW_VOICE_CH")
	)

	cfg, err := config.ReadConfig("config.yaml",
		config.WithServerID(serverID),
		config.WithMainChannelID(mainChannelID),
		config.WithVoiceChannelID(voiceChannelID),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	bot := stewart.New(cfg, token)
	if err := bot.Run(); err != nil {
		logrus.Fatal(err)
	}

	jobCh := job.New(bot.Session, cfg).Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc

	close(jobCh)
	bot.Close()
}
