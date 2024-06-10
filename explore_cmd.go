package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SpecificLocationResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", name)

	visited := false
	for _, loc := range globalConfig.VisitedLocations {
		if loc == url {
			visited = true
			break
		}
	}

	// Cached Response
	if visited {
		val, ok := cache.Get(url)

		if ok {
			fmt.Println("--Explore Cached Response--")

			var pokemonsFound SpecificLocationResponse
			err := json.Unmarshal(val, &pokemonsFound)

			if err != nil {
				return err
			}

			fmt.Println("Pokemons found:")

			for _, pokemonEnc := range pokemonsFound.PokemonEncounters {
				fmt.Printf("- %s\n", pokemonEnc.Pokemon.Name)
			}

			return nil
		}
	}

	// Non-Cached Response
	res, err := http.Get(url)

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
	cache.Add(url, body)
	globalConfig.VisitedLocations = append(globalConfig.VisitedLocations, url)

	var pokemonsFound SpecificLocationResponse
	err = json.Unmarshal(body, &pokemonsFound)

	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Pokemons found:")

	for _, pokemonEnc := range pokemonsFound.PokemonEncounters {
		fmt.Printf("- %s\n", pokemonEnc.Pokemon.Name)
	}

	return nil
}
