package game

type GameResult int

const (
	LOST GameResult = iota - 1
	TIE
	WON
)

type GameChoice struct {
	Value rune
}

func (gc GameChoice) String() string {
	switch gc.Value {
	case 'R':
		return "ROCK"
	case 'P':
		return "PAPER"
	case 'S':
		return "SCISSORS"
	default:
		return "UNKNOWN"
	}
}
