package pokeapi

type ResponseBody struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
}

type LocationsJson struct {
	ResponseBody
	Results  []LocationJson `json:"results"`
}

type LocationJson struct {
	ID                   	int    `json:"id"`
	Name                 	string `json:"name"`
	Location 							struct {
		Name 									string `json:"name"`
		URL  									string `json:"url"`
	} 										`json:"location"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				Chance          int   `json:"chance"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
