package pokeapi

import (
	"encoding/json"
	"errors"

	"github.com/andersjbe/pokedex-cli/internal/pokecache"
)

func GetPokemonByName(name string, cache pokecache.Cache) (Pokemon, error) {
	body, httpError := FetchUrl("https://pokeapi.co/api/v2/pokemon/" + name, &cache)
	if httpError != nil {
		return Pokemon{}, httpError
	}

	pokemon := Pokemon{}
	jsonError := json.Unmarshal(body, &pokemon)
	if jsonError != nil {
		return Pokemon{}, errors.New("Error parsing pokemon JSON")
	}

	return pokemon, nil
}
