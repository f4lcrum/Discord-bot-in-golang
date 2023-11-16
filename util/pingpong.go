package util

import (
	"github.com/bwmarrin/discordgo"
)

func PingPong(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
