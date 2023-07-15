package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GetLocations(url string) ([]string, string, string, error) {
	res, httpErr := http.Get(url)
	if res.StatusCode >= 400 {
		return nil, "", "", errors.New("Encountered error while fetching resource")
	}
	if httpErr != nil {
		return nil, "", "", httpErr
	}

	body, readErr := io.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		return nil, "", "", readErr
	}

	location := LocationJson {}
	jsonErr := json.Unmarshal(body, &location)
	if jsonErr != nil {
		return nil, "", "", jsonErr
	}

	results := make([]string, 20)
	for _, loc := range location.Results {
		results = append(results, loc.Name)
	}

	previous := ""
	if s, ok := location.Previous.(string); ok {
		previous = s
	}
	return results, location.Next, previous, nil
}
