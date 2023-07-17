package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andersjbe/pokedex-cli/internal/pokeapi"
	"github.com/andersjbe/pokedex-cli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	mux := &sync.Mutex{}

	commands := getCommands()
	cache := pokecache.NewCache(time.Minute * 5, mux)
	conf := make(map[string]config)
	conf["locations"] = config{
		next: "https://pokeapi.co/api/v2/location-area/?limit=20",
		previous: "",
	}
	pokedex := make(map[string]pokeapi.Pokemon)

	ctx := context {
		commands: &commands,
		cache: &cache,
		pages: &conf,
		pokedex: &pokedex,
	}

	fmt.Print("pokedex > ")
	for ; scanner.Scan(); fmt.Print("pokedex > ") {
		input := scanner.Text()
		inputs := strings.Split(input, " ")
		args := inputs[1:]
		command, ok := commands[inputs[0]]
		if !ok  {
			fmt.Fprintln(os.Stderr, "command not recognized")
			continue
		}
		err := command.callback(&ctx, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

type context struct {
	commands *map[string]cliCommand
	cache *pokecache.Cache
	pages *map[string]config
	pokedex *map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name 					string
	description 	string
	callback 			func(*context, []string) error
}

type config	struct {
	next 					string
	previous 			string
}

func getCommands() map[string]cliCommand  {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "Output all available commands",
			callback: helpCommand,
		},
		"exit": {
			name: "exit",
			description: "Exit the program",
			callback: exitCommand,
		},
		"map": {
			name: "map",
			description: "Print the next 20 Pokemon game locations",
			callback: mapCommand,
		},
		"mapb": {
			name: "mapb",
			description: "Print the last 20 Pokemon game locations",
			callback: mapBackCommand,
		},
		"explore": {
			name: "explore <location>",
			description: "List all the pokemon found in a location",
			callback: exploreCommand,
		},
		"catch": {
			name: "catch <pokemon>",
			description: "Attempt to catch a pokemon. Success rate is based on pokemon's base experience",
			callback: catchCommand,
		},
		"identify": {
			name: "identify <pokemon>",
			description: "View the details of a caught pokemon",
			callback: identifyCommand,
		},
		"pokedex": {
			name: "pokedex",
			description: "List the pokemon recorded in the pokedex",
			callback: pokedexCommand,
		},
	}
}

func helpCommand(ctx *context, _ []string)  error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range *ctx.commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func exitCommand(_ *context, _ []string) error {
	os.Exit(0)
	return nil
}

func mapCommand(ctx *context, _ []string) error {
	conf := *ctx.pages

	locations, next, previous, err := pokeapi.GetLocations(conf["locations"].next, ctx.cache)
	if err != nil {
		return err
	}

	conf["locations"] = config{
		next: next,
		previous: previous,
	}
	*ctx.pages = conf

	fmt.Println("Locations:")
	for i:=0; i<len(locations); i++ {
		fmt.Println(" - " + locations[i])
	}

	return nil
}

func mapBackCommand(ctx *context, _ []string) error {
	conf := *ctx.pages

	if conf["locations"].previous == "" {
		return errors.New("No previous page")
	}
	locations, next, previous, err := pokeapi.GetLocations(conf["locations"].previous, ctx.cache)
	if err != nil {
		return err
	}

	conf["locations"] = config{
		next: next,
		previous: previous,
	}
	*ctx.pages = conf

	fmt.Println("Locations:")
	for i:=0; i<len(locations); i++ {
		fmt.Println(" - " + locations[i])
	}

	return nil
}

func exploreCommand(ctx *context, args []string) error {
	if len(args) == 0 {
		return errors.New("Please enter a valid location to explore")
	}
	location, err := pokeapi.GetLocationByName("https://pokeapi.co/api/v2/location-area/" + args[0], ctx.cache)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range location.PokemonEncounters {
		fmt.Println(" - " + pokemon.Pokemon.Name)
	}

	return nil
}

func catchCommand(ctx *context, args []string) error {
	if len(args) == 0 {
		return errors.New("Please enter a pokemon to catch")
	}

	pokemon, err := pokeapi.GetPokemonByName(args[0], *ctx.cache)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a pokeball at %s...\n", pokemon.Name)

	// Generate a random number that represents the success of the catch
	r := rand.New(rand.NewSource(time.Now().Unix()))
	catchRoll := r.Intn(637)

	if catchRoll >= pokemon.BaseExperience {
		pokedex := *ctx.pokedex
		pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else if catchRoll < pokemon.BaseExperience && catchRoll > pokemon.BaseExperience + 100 {
		fmt.Println("Shoot! Almost had it")
	} else {
		fmt.Printf("%s broke free!\n", pokemon.Name)
	}

	return nil
}

func identifyCommand(ctx *context, args []string) error {
	if len(args) == 0 {
		return errors.New("Please enter a pokemon to identify")
	}

	pokedex := *ctx.pokedex
	pokemon, found := pokedex[args[0]]
	if !found {
		return errors.New(fmt.Sprintf("Pokemon %s not recognized", args[0]))
	}

	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("\t- %s\n", typ.Type.Name)
	}
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	return nil
}

func pokedexCommand(ctx *context, _ []string) error {
	println("Your Pokedex:")

	for _, pokemon := range *ctx.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}
