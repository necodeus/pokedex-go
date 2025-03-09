package commands

import (
	"fmt"
	"os"

	config "github.com/necodeus/pokedex-go/config"
)

func CommandExit(config *config.Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
