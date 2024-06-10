package main

import (
	"github.com/LucasCoppola/pokedex-cli/internal/pokecache"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
}

type Config struct {
	Next             string
	Previous         *string
	VisitedLocations []string
}

var globalConfig = Config{Next: "https://pokeapi.co/api/v2/location-area", Previous: nil}
var cache = pokecache.NewCache(5 * time.Minute)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List of pokemons caught",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Inspect pokemon",
			callback:    commandInspect,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations in the pokemon world",
			callback:    commandMapBack,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
