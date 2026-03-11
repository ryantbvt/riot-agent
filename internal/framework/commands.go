package framework

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/ryantbvt/riot-agent/internal/riot/lol"
)

type Command struct {
	Name        string
	Description string
	Execute     func(h *Handler, s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

// summonerReviewCommand handles the Discord interaction for summoner review
func summonerReviewCommand(h *Handler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Validate args
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !summoner-review <summonerName#tagLine> [region]")
		return
	}

	riotID := args[0]

	// Default to americas if region is not passed
	region := ""
	if len(args) == 2 {
		region = args[1]
	}
	region, err := validRegion(region)
	if err != nil {
		log.Println("Invalid region")
	}

	result, err := lol.GetSummonerReview(h.RiotClient, riotID, region)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Send result
	s.ChannelMessageSend(m.ChannelID, result)
}

// Validate if the region passed is valid
func validRegion(region string) (string, error) {
	validRegions := []string{
		"americas",
		"europe",
		"asia",
	}

	if region == "" {
		region = "americas"
	}

	if !slices.Contains(validRegions, region) {
		errMsg := strings.Join(validRegions, ", ")
		return "", fmt.Errorf("Invalid region. Region must be: %s", errMsg)
	}

	return region, nil
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"summoner-review": {
			Name:        "summoner-review",
			Description: "Reviews players league of legends match history and rank",
			Execute:     summonerReviewCommand,
		},
	}
}
