package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
)

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

var pokedexMap = make(map[string]Pokemon)

func commandCatch(args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)

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

	var pokemonFound Pokemon
	err = json.Unmarshal(body, &pokemonFound)

	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonFound.Name)

	baseExperience := pokemonFound.BaseExperience

	// Arceus seems to be the most difficult to catch and
	// his base_exp is 324 based on the api
	catchThreshold := 330
	catchChance := catchThreshold - baseExperience

	if catchChance < 0 {
		catchChance = 0
	}

	percentageChange := float64(catchChance) / float64(catchThreshold) * 100
	isCaught := rand.IntN(catchThreshold) < catchChance

	fmt.Printf("You have a %.0f%% chance of catching it.\n", percentageChange)

	if isCaught {
		_, ok := pokedexMap[pokemonFound.Name]

		if ok {
			fmt.Printf("\n%s has already been caught!\n", pokemonFound.Name)
		} else {
			fmt.Printf("\n%s was caught!\n", pokemonFound.Name)
			fmt.Println("You may now inspect it with the inspect command.")
			pokedexMap[pokemonFound.Name] = pokemonFound
		}
	} else {
		fmt.Printf("\n%s escaped!\n", pokemonFound.Name)
	}

	return nil
}
