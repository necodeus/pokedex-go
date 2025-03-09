package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/necodeus/pokedex-go/internal/pokecache"
)

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

func GetPokemon(pokemon string, cache *pokecache.Cache) (*PokemonResponse, error) {
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
