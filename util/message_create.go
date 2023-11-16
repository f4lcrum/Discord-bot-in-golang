package util

import (
	"github.com/bwmarrin/discordgo"
)

// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	PingPong(s, m)
}
