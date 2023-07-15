package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/andersjbe/pokedex-cli/internal/pokecache"
)

func GetLocations(url string, cache pokecache.Cache) ([]string, string, string, error) {
	if _, cached := cache.Get(url); !cached {
		fetched, fetchErr := fetchLocations(url)
		if fetchErr != nil {
			return nil, "", "", fetchErr
		}
		cache.Add(url, fetched)
	}

	body, _ := cache.Get(url)

	location, jsonErr := parseLocations(body)
	if jsonErr != nil {
		return nil, "", "", jsonErr
	}

	results := make([]string, 0, 20)
	for _, loc := range location.Results {
		results = append(results, loc.Name)
	}

	previous := ""
	if s, ok := location.Previous.(string); ok {
		previous = s
	}
	return results, location.Next, previous, nil
}

func fetchLocations(url string) ([]byte, error) {
	res, httpErr := http.Get(url)
	if res.StatusCode >= 400 {
		return nil, errors.New("Encountered error while fetching resource")
	}
	if httpErr != nil {
		return nil, httpErr
	}

	body, readErr := io.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

func parseLocations(body []byte) (LocationJson, error) {
	location := LocationJson {}
	jsonErr := json.Unmarshal(body, &location)
	if jsonErr != nil {
		return location, jsonErr
	}

	return location, nil
}
