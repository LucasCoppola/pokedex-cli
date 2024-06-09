package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}

		if len(words) > 1 {
			args = words[1:]
		}

		commandMap := getCommands()
		command, ok := commandMap[commandName]

		if !ok {
			fmt.Println("Sorry that command doesn't exist, try running `help`")
			continue
		} else {
			err := command.callback(args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

	}

}
