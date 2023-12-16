package moderation

import (
	"bufio"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

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

func GetRule(swearWordsPath string, channelLogId string) discordgo.AutoModerationRule {
	keywords, err := readLines(swearWordsPath)

	if err != nil {
		log.Fatal("error loading swear words,", err)
	}
	enabled := true
	rule := discordgo.AutoModerationRule{
		Name:        "Bad words auto mod",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: keywords,
		},

		Enabled: &(enabled),
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionBlockMessage},
			{Type: discordgo.AutoModerationRuleActionTimeout, Metadata: &discordgo.AutoModerationActionMetadata{Duration: 30}},
			{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{
				ChannelID: channelLogId,
			}},
		},
	}

	return rule
}
