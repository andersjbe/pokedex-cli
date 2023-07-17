package pokeapi

import (
	"encoding/json"
	"errors"

	"github.com/andersjbe/pokedex-cli/internal/pokecache"
)

func GetLocations(url string, cache *pokecache.Cache) ([]string, string, string, error) {
	body, fetchErr := FetchUrl(url, cache)
	if fetchErr != nil {
		return nil, "", "", fetchErr
	}

	locations := LocationsJson {}
	jsonErr := json.Unmarshal(body, &locations)
	if jsonErr != nil {
		return nil, "", "", jsonErr
	}

	results := make([]string, 0, 20)
	for _, loc := range locations.Results {
		results = append(results, loc.Name)
	}

	return results, locations.Next, locations.Previous, nil
}

func GetLocationByName(url string, cache *pokecache.Cache) (LocationJson, error) {
	body, fetchErr := FetchUrl(url, cache)
	if fetchErr != nil {
		return LocationJson{}, errors.New("Location not found")
	}

	location := LocationJson {}
	jsonErr := json.Unmarshal(body, &location)
	if jsonErr != nil {
		return LocationJson{}, jsonErr
	}

	return location, nil

}
