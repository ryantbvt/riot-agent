package framework

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/junioryono/Riot-API-Golang/apiclient"
)

const (
	Prefix = "!"
)

type Handler struct {
	RiotClient apiclient.Client
}

func NewHandler(riotClient apiclient.Client) *Handler {
	return &Handler{
		RiotClient: riotClient,
	}
}

func (h *Handler) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Validate message is not itself
	if m.Author.Bot {
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

	// Validate the command exist
	commands := GetCommands()
	cmd, exists := commands[cmdName]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
		return
	}

	cmd.Execute(h, s, m, args)

}
