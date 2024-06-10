package main

import (
	"errors"
	"fmt"
)

func commandInspect(args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	pokemon, ok := pokedexMap[name]

	if !ok {
		fmt.Printf("%s has not been caught yet\n", name)
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, v := range pokemon.Stats {
			fmt.Printf(" - %s: %d\n", v.Stat.Name, v.BaseStat)
		}
		fmt.Println("Types:")
		for _, v := range pokemon.Types {
			fmt.Printf(" - %s\n", v.Type.Name)
		}
	}

	return nil
}
