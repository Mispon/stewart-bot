package main

import (
	"errors"
	"flag"
	"github.com/mispon/stewart-bot/internal/job"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	stewart "github.com/mispon/stewart-bot/internal/bot"
	"github.com/mispon/stewart-bot/internal/config"
)

var (
	debug          bool
	token          string
	serverID       string
	mainChannelID  string
	voiceChannelID string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "--debug=true")
	flag.StringVar(&token, "token", "", "--token=my_bot_token")
	flag.StringVar(&serverID, "server_id", "", "--server_id=my_server_id")
	flag.StringVar(&mainChannelID, "main_ch", "", "--main_ch=my_main_channel_id")
	flag.StringVar(&voiceChannelID, "voice_ch", "", "--voice_ch=my_voice_channel_id")

	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	logrus.Infoln("Initialize bot...")

	if err := validateFlags(); err != nil {
		logrus.Fatal(err)
	}

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

func validateFlags() error {
	if len(token) == 0 {
		return errors.New("token is empty")
	}

	if len(serverID) == 0 {
		return errors.New("server id is empty")
	}

	if len(mainChannelID) == 0 {
		return errors.New("main channel id is empty")
	}

	if len(voiceChannelID) == 0 {
		return errors.New("voice channel id is empty")
	}

	return nil
}
