package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LucasCoppola/pokedex-cli/internal/pokecache"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type Config struct {
	Next     string
	Previous *string
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

var globalConfig = Config{Next: "https://pokeapi.co/api/v2/location", Previous: nil}
var cache = pokecache.NewCache(5 * time.Minute)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
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
	}
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMap() error {
	val, ok := cache.Get(globalConfig.Next)

	// Cached Response
	if ok {
		fmt.Println("--Map Cached Response--")
		var locations LocationResponse
		err := json.Unmarshal(val, &locations)

		if err != nil {
			return err
		}

		globalConfig.Previous = locations.Previous
		globalConfig.Next = locations.Next

		for _, location := range locations.Results {
			fmt.Println(location.Name)
		}

		return nil
	}

	// Non-Cached Response
	res, err := http.Get(globalConfig.Next)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	// Add to Cache
	cache.Add(globalConfig.Next, body)

	var locations LocationResponse
	err = json.Unmarshal(body, &locations)

	if err != nil {
		return err
	}

	globalConfig.Previous = locations.Previous
	globalConfig.Next = locations.Next

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack() error {
	if globalConfig.Previous == nil {
		return errors.New("You're on the first page")
	}

	val, ok := cache.Get(*globalConfig.Previous)

	// Cached Response
	if ok {
		fmt.Println("--Map Back Cached Response--")
		var locations LocationResponse
		err := json.Unmarshal(val, &locations)

		if err != nil {
			return err
		}

		globalConfig.Previous = locations.Previous
		globalConfig.Next = locations.Next

		for _, location := range locations.Results {
			fmt.Println(location.Name)
		}

		return nil
	}

	// Non-Cached Response
	res, err := http.Get(*globalConfig.Previous)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	// Add to Cache
	cache.Add(*globalConfig.Previous, body)

	var locations LocationResponse
	err = json.Unmarshal(body, &locations)

	if err != nil {
		return err
	}

	globalConfig.Previous = locations.Previous
	globalConfig.Next = locations.Next

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}
