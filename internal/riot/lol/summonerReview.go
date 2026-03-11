package lol

import (
	"fmt"
	"log"
	"strings"

	"github.com/junioryono/Riot-API-Golang/apiclient"
	"github.com/junioryono/Riot-API-Golang/constants/continent"
	"github.com/ryantbvt/riot-agent/internal/riot"
)

type MatchReview struct {
	Champion string  `json:"champion"`
	Role     string  `json:"role"`
	Win      bool    `json:"win"`
	Duration float64 `json:"duration_minutes"`

	GoldPerMin float64 `json:"gold_per_min"`
	CsPerMin   float64 `json:"cs_per_min"`

	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`

	DamagePerMin      float64 `json:"damage_per_min"`
	DamageShare       float64 `json:"damage_share"`
	KillParticipation float64 `json:"kill_participation"`

	GoldDiff15        int     `json:"gold_diff_15"`
	XpDiff15          int     `json:"xp_diff_15"`
	VisionScorePerMin float64 `json:"vision_score_per_min"`
}

// GetSummonerReview analyzes summoner's recent matches and returns the review data using the Riot-API-Golang client.
func GetSummonerReview(client apiclient.Client, riotID string, region string) (string, error) {
	// Parse riotID format: summonerName#tagLine
	parts := strings.SplitN(riotID, "#", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: expected summonerName#tagLine")
	}
	summonerName, tagLine := parts[0], parts[1]

	cont, err := riot.RegionToContinent(region)
	if err != nil {
		return "", fmt.Errorf("region: %w", err)
	}

	// Get account (puuid) via Riot-API-Golang
	account, err := client.GetAccountByRiotID(cont, summonerName, tagLine)
	if err != nil {
		return "", fmt.Errorf("failed to fetch summoner: %w", err)
	}
	puuid := account.Puuid

	// Get match list (built-in rate limiting)
	matchlist, err := client.GetMatchlist(cont, puuid, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get player match history: %w", err)
	}
	if matchlist == nil || len(*matchlist) == 0 {
		return "", fmt.Errorf("no matches found for summoner")
	}

	// Get match details and build review data
	reviewData, err := getMatchReviews(client, cont, *matchlist, puuid)
	if err != nil {
		return "", fmt.Errorf("failed to get match info: %w", err)
	}

	// TODO: Send to LLM to review the matches
	// Compile and send response back
	result := reviewData[0].Champion

	log.Printf("Successfully reviewed player: %s", summonerName)
	return result, nil
}

func getMatchReviews(client apiclient.Client, cont continent.Continent, matchIDs []string, puuid string) ([]MatchReview, error) {
	var reviewData []MatchReview
	log.Printf("Reviewing %d matches", len(matchIDs))

	for _, matchID := range matchIDs {
		match, err := client.GetMatch(cont, matchID)
		if err != nil {
			continue // skip failed matches
		}
		if match == nil || match.Info.Participants == nil {
			continue
		}

		info := &match.Info
		durationMin := float64(info.GameDuration) / 60.0
		if durationMin <= 0 {
			continue
		}

		var player *apiclient.MatchInfoParticipant
		var totalTeamDamage int32

		for i := range info.Participants {
			p := &info.Participants[i]
			if p.TeamID == 100 {
				totalTeamDamage += p.TotalDamageDealtToChampions
			}
			if p.SummonerPuuid == puuid {
				player = p
			}
		}
		if player == nil {
			continue
		}

		cs := int32(player.TotalMinionsKilled) + int32(player.NeutralMinionsKilled)
		damage := player.TotalDamageDealtToChampions
		gold := player.GoldEarned
		damageShare := 0.0
		if totalTeamDamage > 0 {
			damageShare = float64(damage) / float64(totalTeamDamage)
		}

		kp := 0.0
		visionPerMin := 0.0
		if player.Challenges != nil {
			kp = player.Challenges.KillParticipation
			visionPerMin = float64(player.Challenges.VisionScorePerMinute)
		}
		// goldDiffAt15 / xpDiffAt15 not in Riot-API-Golang Challenges struct; keep as 0
		goldDiff15 := 0
		xpDiff15 := 0

		review := MatchReview{
			Champion: player.ChampionName,
			Role:     player.TeamPosition,
			Win:      player.Win,
			Duration: durationMin,

			GoldPerMin: float64(gold) / durationMin,
			CsPerMin:   float64(cs) / durationMin,

			Kills:   int(player.Kills),
			Deaths:  int(player.Deaths),
			Assists: int(player.Assists),

			DamagePerMin:      float64(damage) / durationMin,
			DamageShare:       damageShare,
			KillParticipation: kp,

			GoldDiff15:        goldDiff15,
			XpDiff15:          xpDiff15,
			VisionScorePerMin: visionPerMin,
		}
		reviewData = append(reviewData, review)
	}

	return reviewData, nil
}
