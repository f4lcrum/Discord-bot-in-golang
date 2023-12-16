package main

import (
	"discord_bot/branik"
	"discord_bot/commands"
	"discord_bot/moderation"
	"discord_bot/reactions"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Token          string `env:"TOKEN,required"`
	Guild          string `env:"GUILD,required"`
	ChannelLog     string `env:"CHANNEL_LOG,required"`
	SwearWordsPath string `env:"SWEAR_WORDS_PATH,required"`
	ApiKey         string `env:"API_KEY,required"`
}

var (
	cfg = Config{}
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
		return
	}

	dg.AddHandler(reactions.SaveByReaction)
	dg.AddHandler(branik.Branig)

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.Identify.Intents |= discordgo.IntentGuildMessageReactions

	rule := moderation.GetRule(cfg.SwearWordsPath, cfg.ChannelLog)
	moderation.AddRule(dg, rule, cfg.Guild)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	cmds := commands.AddCommands(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	commands.RemoveCommands(dg, cmds)
	moderation.RemoveRule(dg, rule.ID, cfg.Guild)
	err = dg.Close()
	if err != nil {
		return
	}
}
