package commands

import config "github.com/necodeus/pokedex-go/config"

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*config.Config, []string) error
}

var Commands map[string]CliCommand
