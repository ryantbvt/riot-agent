package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	riotUrl = "https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s"
)

// Get the PUUID using summonerName#tagLine
func (r *Client) GetRiotId(summonerName string, tagLine string, region string) (string, error) {
	url := fmt.Sprintf(riotUrl, region, summonerName, tagLine)

	// Create HTTP request with API key header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Riot-Token", r.apiKey)

	// Execute request
	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		PUUID string `json:"puuid"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Successfully fetched PUUID for: %s", summonerName)
	return result.PUUID, nil
}
