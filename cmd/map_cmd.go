package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

func commandMap(args ...string) error {
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

func commandMapBack(args ...string) error {
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
