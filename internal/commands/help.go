package commands

import (
	"fmt"

	config "github.com/necodeus/pokedex-go/config"
)

func CommandHelp(config *config.Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n")
	for _, command := range Commands {
		fmt.Printf(" - %s: %s\n", command.Name, command.Description)
	}
	return nil
}
