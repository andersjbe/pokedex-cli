package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

	ctx := context {
		commands: &commands,
		cache: &cache,
		pages: &conf,
	}

	fmt.Print("pokedex > ")
	for ; scanner.Scan(); fmt.Print("pokedex > ") {
		input := scanner.Text()
		command, ok := commands[input]
		if !ok  {
			fmt.Fprintln(os.Stderr, "command not recognized")
			continue
		}
		err := command.callback(&ctx)
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
}

type cliCommand struct {
	name 					string
	description 	string
	callback 			func(*context) error
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
	}
}

func helpCommand(ctx *context)  error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range *ctx.commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func exitCommand(_ *context) error {
	os.Exit(0)
	return nil
}

func mapCommand(ctx *context) error {
	conf := *ctx.pages

	locations, next, previous, err := pokeapi.GetLocations(conf["locations"].next, *ctx.cache)
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

func mapBackCommand(ctx *context) error {
	conf := *ctx.pages

	if conf["locations"].previous == "" {
		return errors.New("No previous page")
	}
	locations, next, previous, err := pokeapi.GetLocations(conf["locations"].previous, *ctx.cache)
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
