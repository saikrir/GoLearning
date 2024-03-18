package main

import (
	"RPS/pkg/game"
	"fmt"
)

func printGameMenu() {
	fmt.Println("Welcome to Rock Paper Scisccors Game")
	fmt.Println("Select one of the choices ")
	fmt.Println("Rock --> R ")
	fmt.Println("Paper --> P ")
	fmt.Println("Scissors --> S ")
}

func getUserChoice() rune {
	var choice rune
	for {
		fmt.Printf("\nPlease enter your choice: ")
		_, err := fmt.Scanf("%c\n", &choice)

		if err != nil {
			fmt.Println("\nInvalid input, please try again ", err)
			// If int read fails, nullify rest of the input
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if !game.IsValidGameMove(choice) {
			fmt.Printf("\nInvalid choice, R, P, S are only valid choices, You entered %c \n", choice)
			continue
		}
		break
	}

	return choice
}

func main() {
	printGameMenu()
	fmt.Printf("Press Ctrl + C to stop")
	game.RegisterExitHandler()
	for {
		userMove := game.GameChoice{Value: getUserChoice()}
		fmt.Println(game.GetResult(userMove))
	}
}
