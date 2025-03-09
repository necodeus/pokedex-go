package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text) // remove leading and trailing spaces
	text = strings.ToLower(text)   // convert to lowercase

	return strings.Fields(text) // split the text into words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getLocationAreas(url string) (*LocationAreaResponse, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("unexpected status code: %d", res.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return &data, nil
}

func commandMap(config *config) error {
	data, err := getLocationAreas(config.Next)
	if err != nil {
		return err
	}

	config.Next = data.Next
	config.Previous = data.Previous

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapBack(config *config) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
	}

	data, err := getLocationAreas(config.Previous)
	if err != nil {
		return err
	}

	config.Next = data.Next
	config.Previous = data.Previous

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

type config struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand

func main() {
	// define the commands
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of location areas",
			callback:    commandMapBack,
		},
	}

	// wait for user input
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{}

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

		// check if the command exists
		command, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		// execute the command
		err := command.callback(&cfg)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
