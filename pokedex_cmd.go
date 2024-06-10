package main

import "fmt"

func commandPokedex(args ...string) error {
	fmt.Println("Your pokedex:")
	for _, pokemon := range pokedexMap {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}
