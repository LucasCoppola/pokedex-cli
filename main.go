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

		commandMap := getCommands()
		command, ok := commandMap[words[0]]

		if !ok {
			fmt.Println("Sorry that command doesn't exist, try running `help`")
			continue
		} else {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

	}

}
