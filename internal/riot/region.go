package riot

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/constants/continent"
)

// RegionToContinent maps a region string (americas, europe, asia) to the Riot API continent.
func RegionToContinent(region string) (continent.Continent, error) {
	switch region {
	case "americas":
		return continent.AMERICAS, nil
	case "europe":
		return continent.EUROPE, nil
	case "asia":
		return continent.ASIA, nil
	default:
		return "", fmt.Errorf("invalid region: %s", region)
	}
}
