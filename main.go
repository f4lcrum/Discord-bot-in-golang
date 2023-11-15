package main

import (
	"bufio"
	"discord_bot/util"
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
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing swear words file")
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := Config{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
		return
	}

	dg.AddHandler(util.MessageCreate)

	enabled := true

	keywords, err := readLines(cfg.SwearWordsPath)

	if err != nil {
		fmt.Println("error loading swear words,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	rule, err := dg.AutoModerationRuleCreate(cfg.Guild, &discordgo.AutoModerationRule{
		Name:        "Auto Moderation",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: keywords,
		},

		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionBlockMessage},
			{Type: discordgo.AutoModerationRuleActionTimeout, Metadata: &discordgo.AutoModerationActionMetadata{Duration: 30}},
			{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{
				ChannelID: cfg.ChannelLog,
			}},
		},
	})

	util.AddRule(dg, rule, cfg.Guild)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	commands := util.AddCommands(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	util.RemoveCommands(dg, commands)
	util.RemoveRule(dg, rule.ID, cfg.Guild)
	err = dg.Close()
	if err != nil {
		return
	}
}
