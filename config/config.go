package config

import (
	"github.com/necodeus/pokedex-go/internal/pokeapi"
	"github.com/necodeus/pokedex-go/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokemon  map[string]pokeapi.PokemonResponse // for now we use PokemonResponse but i don't think we need all the data
}
