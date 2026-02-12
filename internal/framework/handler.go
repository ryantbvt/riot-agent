package framework

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ryantbvt/riot-agent/internal/riot"
	"github.com/ryantbvt/riot-agent/internal/riot/lol"
)

const (
	Prefix = "!"
)

type Handler struct {
	Riot      *riot.Client
	LolClient *lol.LolClient
}

func NewHandler(riotClient *riot.Client, LolClient *lol.LolClient) *Handler {
	return &Handler{
		Riot:      riotClient,
		LolClient: LolClient,
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
