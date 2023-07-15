package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/andersjbe/pokedex-cli/internal/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	getCommands := createCommands()
	commands := *getCommands()

	fmt.Print("pokedex > ")
	for ; scanner.Scan(); fmt.Print("pokedex > ") {
		input := scanner.Text()
		command, ok := commands[input]
		if !ok  {
			fmt.Fprintln(os.Stderr, "command not recognized")
			continue
		}
		err := command.callback(&commands)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

type cliCommand struct {
	name 					string
	description 	string
	callback 			func(*map[string]cliCommand) error
	config				*config
}

type config	struct {
	next 					string
	previous 			string
}

func createCommands() func() *map[string]cliCommand  {
	locationConfig := config{
		next: "https://pokeapi.co/api/v2/location-area/?limit=20",
		previous: "",
	}

	commands := map[string]cliCommand{
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
			config: &locationConfig,
		},
		"mapb": {
			name: "mapb",
			description: "Print the last 20 Pokemon game locations",
			callback: mapBackCommand,
			config: &locationConfig,
		},
	}

	return func() *map[string]cliCommand {
		return &commands
	}
}

func helpCommand(commandRefs *map[string]cliCommand) error {
	commands := *commandRefs

	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func exitCommand(_  *map[string]cliCommand) error {
	os.Exit(0)
	return nil
}

func mapCommand(commandRefs *map[string]cliCommand) error {
	commands := *commandRefs

	locations, next, previous, err := pokeapi.GetLocations(commands["map"].config.next)

	if err != nil {
		return err
	}

	commands["map"].config.next = next
	commands["map"].config.previous = previous

	for _, location := range locations {
		fmt.Println(location)
	}

	return nil
}

func mapBackCommand(commandRefs *map[string]cliCommand) error {
	commands := *commandRefs

	if commands["map"].config.previous == "" {
		return errors.New("No previous page")
	}
	locations, next, previous, err := pokeapi.GetLocations(commands["map"].config.previous)

	if err != nil {
		return err
	}

	commands["map"].config.next = next
	commands["map"].config.previous = previous

	for _, location := range locations {
		fmt.Println(location)
	}

	return nil
}
