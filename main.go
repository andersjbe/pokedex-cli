package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	fmt.Print("pokedex > ")
	for ; scanner.Scan(); fmt.Print("pokedex > ") {
		input := scanner.Text()
		command, ok := commands[input]
		if !ok  {
			fmt.Fprintln(os.Stderr, "command not recognized")
			continue
		}
		command.callback()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

type cliCommand struct {
	name 					string
	description 	string
	callback 			func() error
}

func getCommands() map[string]cliCommand {
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
	}
}

func helpCommand() error {
	commands := getCommands()

	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func exitCommand() error {
	os.Exit(0)
	return nil
}
