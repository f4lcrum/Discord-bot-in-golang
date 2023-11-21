package branik

import (
	"discord_bot/branik/branikParser"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Branig(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	branikId := os.Getenv("CHANNEL_BRANIK")

	if m.Author.ID == s.State.User.ID || m.ChannelID != branikId {
		return
	}
	s.ChannelMessageSend(m.ChannelID, branikParser.Parse(m.Content))

}
