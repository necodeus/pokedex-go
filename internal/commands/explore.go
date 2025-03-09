package commands

import (
	"fmt"

	config "github.com/necodeus/pokedex-go/config"
	"github.com/necodeus/pokedex-go/internal/pokeapi"
)

func CommandExplore(config *config.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: explore <location_area>")
		return nil
	}

	locationArea := args[0]

	fmt.Printf("Exploring %s...\n", locationArea)
	data, err := pokeapi.GetLocationArea(locationArea, config.Cache)
	if err != nil {
		fmt.Println("Exploration failed:", err)
		return nil
	}

	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this location area")
		return nil
	}

	fmt.Println("Pokémon found:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
