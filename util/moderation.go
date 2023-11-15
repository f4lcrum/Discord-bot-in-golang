package util

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func AutoModFuncExecution(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
	channel, err := s.UserChannelCreate(e.UserID)

	if err != nil {
		fmt.Println("error creating channel:", err)
		_, err := s.ChannelMessageSend(
			e.ChannelID,
			"Something went wrong while sending the DM!",
		)
		if err != nil {
			return
		}
		return
	}

	if e.Action.Type == discordgo.AutoModerationRuleActionSendAlertMessage {

		message := fmt.Sprintf("Content you wrote: \"%s\" \n is forbidden! Repeated violation of the rules will lead to a ban \n", e.Content)

		_, err = s.ChannelMessageSend(channel.ID, message)
		if err != nil {
			fmt.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				e.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		}
	}
}

func RemoveRule(session *discordgo.Session, ruleId string, guildId string) {
	session.AutoModerationRuleDelete(guildId, ruleId)
	log.Println("Rules removed")
}

func AddRule(session *discordgo.Session, rule *discordgo.AutoModerationRule, guildId string) {
	log.Println("Adding auto moderation rules...")
	session.Identify.Intents |= discordgo.IntentAutoModerationExecution
	session.Identify.Intents |= discordgo.IntentMessageContent
	session.AddHandler(AutoModFuncExecution)

	log.Println("Adding auto moderation rules complete")
}
