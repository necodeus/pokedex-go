package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/necodeus/pokedex-go/internal/pokecache"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text) // remove leading and trailing spaces
	text = strings.ToLower(text)   // convert to lowercase

	return strings.Fields(text) // split the text into words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func commandExit(config *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getLocationAreas(url string, cache *pokecache.Cache) (*LocationAreasResponse, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	if cachedData, found := cache.Get(url); found {
		var cachedResponse LocationAreasResponse
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return nil, err
		}
		return &cachedResponse, nil
	}

	fmt.Println("=======================================================")
	fmt.Println("Sending request to", url)
	fmt.Println("=======================================================")

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var data LocationAreasResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// add the response to the cache
	cache.Add(url, body)

	return &data, nil
}

type LocationAreaResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	GameIndex           int    `json:"game_index"`
	EncounterMethodRate []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func getLocationArea(area string, cache *pokecache.Cache) (*LocationAreaResponse, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)

	if cachedData, found := cache.Get(url); found {
		var cachedResponse LocationAreaResponse
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return nil, err
		}
		return &cachedResponse, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	return &data, nil
}

type PokemonResponse struct {
	ID                     int              `json:"id"`
	Name                   string           `json:"name"`
	BaseExperience         int              `json:"base_experience"`
	Height                 int              `json:"height"`
	IsDefault              bool             `json:"is_default"`
	Order                  int              `json:"order"`
	Weight                 int              `json:"weight"`
	Abilities              []Ability        `json:"abilities"`
	Forms                  []Form           `json:"forms"`
	GameIndices            []GameIndex      `json:"game_indices"`
	HeldItems              []HeldItem       `json:"held_items"`
	LocationAreaEncounters string           `json:"location_area_encounters"`
	Moves                  []Move           `json:"moves"`
	Species                NamedAPIResource `json:"species"`
	Sprites                Sprites          `json:"sprites"`
	Cries                  Cries            `json:"cries"`
	Stats                  []Stat           `json:"stats"`
	Types                  []TypeSlot       `json:"types"`
	PastTypes              []PastType       `json:"past_types"`
}

type Ability struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedAPIResource `json:"ability"`
}

type Form struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GameIndex struct {
	GameIndex int              `json:"game_index"`
	Version   NamedAPIResource `json:"version"`
}

type HeldItem struct {
	Item           NamedAPIResource `json:"item"`
	VersionDetails []VersionDetail  `json:"version_details"`
}

type VersionDetail struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

type Move struct {
	Move                NamedAPIResource     `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	VersionGroup    NamedAPIResource `json:"version_group"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Sprites struct {
	FrontDefault string       `json:"front_default"`
	FrontShiny   string       `json:"front_shiny"`
	Other        OtherSprites `json:"other"`
}

type OtherSprites struct {
	DreamWorld      ImageResource `json:"dream_world"`
	Home            ImageResource `json:"home"`
	OfficialArtwork ImageResource `json:"official-artwork"`
	Showdown        ImageResource `json:"showdown"`
}

type ImageResource struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type Stat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type TypeSlot struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type PastType struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []TypeSlot       `json:"types"`
}

func getPokemon(pokemon string, cache *pokecache.Cache) (*PokemonResponse, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	if cachedData, found := cache.Get(url); found {
		var cachedResponse PokemonResponse
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return nil, err
		}
		return &cachedResponse, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	var data PokemonResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	return &data, nil
}

func commandMap(config *config, args []string) error {
	data, err := getLocationAreas(config.Next, config.Cache)
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

func commandMapBack(config *config, args []string) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
	}

	data, err := getLocationAreas(config.Previous, config.Cache)
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

func commandExplore(config *config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: explore <location_area>")
		return nil
	}

	locationArea := args[0]

	fmt.Printf("Exploring %s...\n", locationArea)
	data, err := getLocationArea(locationArea, config.Cache)
	if err != nil {
		fmt.Println("No encounters found")
		return nil
	}

	if len(data.PokemonEncounters) == 0 {
		return nil
	}

	fmt.Println("Pokemon encounters:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: catch <pokemon>")
		return nil
	}

	pokemon := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	data, err := getPokemon(pokemon, config.Cache)
	if err != nil {
		fmt.Println("Error catching the Pokémon:", err)
		return nil
	}

	// max experience is 608
	// https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_in_Generation_IX
	catchRate := float64(data.BaseExperience) / 608.0
	if catchRate > 1.0 { // we make sure the base experience is not higher than 608
		catchRate = 1.0
	}

	randomNumber := rand.Float64()
	if randomNumber > catchRate {
		fmt.Printf("%s was caught!\n", pokemon)

		config.Pokemon[pokemon] = *data
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}

func commandInspect(config *config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: inspect <pokemon>")
		return nil
	}

	pokemon := args[0]

	data, ok := config.Pokemon[pokemon]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, slot := range data.Types {
		fmt.Printf("  - %s\n", slot.Type.Name)
	}

	return nil
}

func commandPokedex(config *config, args []string) error {
	if len(config.Pokemon) == 0 {
		fmt.Println("You haven't caught any Pokémon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range config.Pokemon {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokemon  map[string]PokemonResponse // for now we use PokemonResponse but i don't think we need all the data
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
		"explore": {
			name:        "explore",
			description: "Displays the Pokémon located in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catches a Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays the details of a caught Pokémon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays names of caught Pokémon",
			callback:    commandPokedex,
		},
	}

	// wait for user input
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{
		Cache:   pokecache.NewCache(10 * time.Second), // create a new cache with a 10 second interval
		Pokemon: make(map[string]PokemonResponse),     // create a map to store the caught Pokémon
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

		// check if the command exists
		command, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		args := words[1:]

		// execute the command
		err := command.callback(&cfg, args)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
