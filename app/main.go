package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"stewart-bot/app/commands"
	"stewart-bot/app/config"
)

var (
	DebugMode  bool
	Token      string
	ConfigPath string
)

func init() {
	flag.BoolVar(&DebugMode, "debug", false, "-debug")
	flag.StringVar(&Token, "token", "", "-token=my_bot_token")
	flag.StringVar(&ConfigPath, "config", "./config.yml", "-config=./config.yml")

	flag.Parse()
}

func main() {
	fmt.Println("Initialize bot...")
	cfgPath, _ := filepath.Abs(ConfigPath)

	if DebugMode {
		fmt.Printf("Config: %s\n", cfgPath)
		fmt.Printf("Token: %s\n", Token)
	}

	config.ReadConfig(cfgPath)

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.AddHandler(commands.OnMessage)

	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Stewart v%s successfully started!", config.GetConfig().Version)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}