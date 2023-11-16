package util

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func SaveByReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	if r.Emoji.Name == "🔖" {
		channel, err := s.UserChannelCreate(r.UserID)
		if err != nil {
			fmt.Println("error creating channel:", err)
			_, err := s.ChannelMessageSend(
				r.ChannelID,
				"Something went wrong while sending the DM!",
			)
			if err != nil {
				return
			}
			return
		}
		msg, _ := s.ChannelMessage(r.ChannelID, r.MessageID)

		message := fmt.Sprintf("------------------------------\nAuthor of message: %s \nContent of message: %s \n", msg.Author.Username, msg.Content)

		for i := 0; i < len(msg.Attachments); i++ {
			message += fmt.Sprintf("Name of file: %s\nLink: %s\n", msg.Attachments[i].Filename, msg.Attachments[i].URL)
		}
		message += "------------------------------"
		_, err = s.ChannelMessageSend(channel.ID, message)
		if err != nil {
			fmt.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				r.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		}
	}
}
