package lol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ryantbvt/riot-agent/internal/riot"
)

const (
	matchHistoryUrl = "https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids"
)

// GetSummonerReview analyzes summoner's recent matches and returns the review data
func GetSummonerReview(client *riot.Client, lolClient *LolClient, riotID string, region string) (string, error) {
	// Parse riotID format: summonerName#tagLine
	parts := strings.SplitN(riotID, "#", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: expected summonerName#tagLine")
	}
	summonerName, tagLine := parts[0], parts[1]

	// Get the summoner puuid
	puuid, err := client.GetRiotId(summonerName, tagLine, region)
	if err != nil {
		return "", fmt.Errorf("failed to fetch summoner: %w", err)
	}

	// Get the match history
	matchIDs, err := lolClient.getPlayerMatchHistory(puuid, region)
	if err != nil {
		return "", fmt.Errorf("failed to get player match history: %w", err)
	}

	// Get dump of the matches

	// Send to LLM to review the matches

	// Compile and send response back
	result := matchIDs[0]

	return result, nil
}

// Get the match history of a player
func (lClient *LolClient) getPlayerMatchHistory(puuid string, region string) ([]string, error) {
	url := fmt.Sprintf(matchHistoryUrl, region, puuid)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Riot-Token", lClient.apiKey)

	// Execute request
	resp, err := lClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var matchIDs []string
	if err := json.NewDecoder(resp.Body).Decode(&matchIDs); err != nil {
		return nil, fmt.Errorf("failed to decode response %w", err)
	}

	return matchIDs, nil
}

func getMatchInfo(matchIDs []string, region string) {
	// pass
}
