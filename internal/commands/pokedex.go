package commands

import (
	"fmt"

	config "github.com/necodeus/pokedex-go/config"
)

func CommandPokedex(config *config.Config, args []string) error {
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
