package commands

import (
	"fmt"
	"math/rand/v2"

	config "github.com/necodeus/pokedex-go/config"
	"github.com/necodeus/pokedex-go/internal/pokeapi"
)

func CommandCatch(config *config.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: catch <pokemon>")
		return nil
	}

	pokemon := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	data, err := pokeapi.GetPokemon(pokemon, config.Cache)
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
