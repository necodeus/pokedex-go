package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/necodeus/pokedex-go/config"
	"github.com/necodeus/pokedex-go/internal/commands"
	"github.com/necodeus/pokedex-go/internal/pokeapi"
	"github.com/necodeus/pokedex-go/internal/pokecache"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text) // remove leading and trailing spaces
	text = strings.ToLower(text)   // convert to lowercase

	return strings.Fields(text) // split the text into words
}

func main() {
	// define the commands
	commands.Commands = map[string]commands.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commands.CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commands.CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays the location areas",
			Callback:    commands.CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous page of location areas",
			Callback:    commands.CommandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Displays the Pokémon located in a location area",
			Callback:    commands.CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catches a Pokémon",
			Callback:    commands.CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Displays the details of a caught Pokémon",
			Callback:    commands.CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Displays names of caught Pokémon",
			Callback:    commands.CommandPokedex,
		},
	}

	// wait for user input
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config.Config{
		Cache:   pokecache.NewCache(10 * time.Second),     // create a new cache with a 10 second interval
		Pokemon: make(map[string]pokeapi.PokemonResponse), // create a map to store the caught Pokémon
	}

	for {
		fmt.Print("Pokedex > ")

		// read the input
		scanner.Scan()
		input := scanner.Text()

		// clean the input
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		// check if the command is defined
		command, ok := commands.Commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		args := words[1:]

		// execute the command
		err := command.Callback(&cfg, args)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
