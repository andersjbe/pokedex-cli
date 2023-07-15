package pokeapi

import (
	"errors"
	"io"
	"net/http"

	"github.com/andersjbe/pokedex-cli/internal/pokecache"
)

func FetchUrl(url string, cache *pokecache.Cache) ([]byte, error) {
	if body, cached := cache.Get(url); cached {
		return body, nil
	}

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

	cache.Add(url, body)
	return body, nil
}
