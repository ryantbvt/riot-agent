package framework

import (
	"github.com/bwmarrin/discordgo"

	"github.com/ryantbvt/riot-agent/internal/riot/lol"
)

type Command struct {
	Name        string
	Description string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

var Commands = map[string]Command{
	"lol-summoner-review": {
		Name:        "summoner-review",
		Description: "Reviews players league of legends match history and rank",
		Execute:     lol.SummonerReview,
	},
}
