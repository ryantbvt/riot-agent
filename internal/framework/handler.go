package framework

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	Prefix = "!"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Validate message is not itself
	if m.Author.Bot {
		log.Println("ignore")
		return
	}

	// Ignore messages w/o the prefix
	if !strings.HasPrefix(m.Content, Prefix) {
		return
	}

	// Remove prefix and split
	content := strings.TrimPrefix(m.Content, Prefix)
	parts := strings.Fields(content)

	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	args := parts[1:]

	cmd, exists := Commands[cmdName]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
		return
	}

	cmd.Execute(s, m, args)

}
