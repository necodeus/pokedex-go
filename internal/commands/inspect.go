package commands

import (
	"fmt"

	config "github.com/necodeus/pokedex-go/config"
)

func CommandInspect(config *config.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: inspect <pokemon>")
		return nil
	}

	pokemon := args[0]

	data, ok := config.Pokemon[pokemon]
	if !ok {
		fmt.Println("You haven't caught this Pokémon yet!")
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
