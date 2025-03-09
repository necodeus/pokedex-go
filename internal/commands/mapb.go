package commands

import (
	"fmt"

	config "github.com/necodeus/pokedex-go/config"
	"github.com/necodeus/pokedex-go/internal/pokeapi"
)

func CommandMapb(config *config.Config, args []string) error {
	if config.Previous == "" {
		fmt.Println("There is no previous page")
	}

	data, err := pokeapi.GetLocationAreas(config.Previous, config.Cache)
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
