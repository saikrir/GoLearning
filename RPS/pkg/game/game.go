package game

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

// P > R
// S > P
// R > S

var validChoices [3]GameChoice = [3]GameChoice{{'R'}, {'P'}, {'S'}}

func (thisChoice GameChoice) Compare(otherMove GameChoice) GameResult {

	if thisChoice.Value == otherMove.Value {
		return TIE
	}

	switch thisChoice.Value {

	case 'R':
		switch otherMove.Value {
		case 'S':
			return WON
		case 'P':
			return LOST
		}

	case 'P':
		switch otherMove.Value {
		case 'R':
			return WON
		case 'S':
			return LOST
		}

	case 'S':
		switch otherMove.Value {
		case 'P':
			return WON
		case 'R':
			return LOST
		}
	}

	return LOST
}

func IsValidGameMove(choice rune) bool {
	for _, validChoice := range validChoices {
		if choice == validChoice.Value {
			return true
		}
	}

	return false
}

func GetResult(userMove GameChoice) string {
	computersMove := MachineMove()
	choicesMade := fmt.Sprintf("You chose %s, Computer chose %s,\nResult: ", userMove, computersMove)

	switch userMove.Compare(computersMove) {
	case TIE:
		return choicesMade + "Game tied"
	case WON:
		return choicesMade + "You won!"
	case LOST:
		return choicesMade + "You lost!"
	default:
		return "Invalid outcome"
	}
}

func MachineMove() GameChoice {
	return validChoices[rand.Intn(2)]
}

func RegisterExitHandler() {

	intteruptChannel := make(chan os.Signal, 2)
	signal.Notify(intteruptChannel, os.Interrupt, syscall.SIGTERM)

	exitHandler := func() {
		<-intteruptChannel
		fmt.Println("\nThanks for playing the Game!")
		os.Exit(0)
	}
	go exitHandler()
}
