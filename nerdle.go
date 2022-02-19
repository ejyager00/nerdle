package main

import (
	"fmt"

	"github.com/ejyager00/nerdle/nerdlegame"
)

const PUZZLE_LENGTH int = 8
const MAX_GUESSES int = 6
const LEADING_ZEROS_ALLOWED bool = false
const NEGATIVE_ZERO_ALLOWED bool = false
const ZERO_REMOVAL_RATE float32 = 1

func main() {
	for i := 0; i < 20; i++ {
		p := nerdlegame.WeightedRandomPuzzle(PUZZLE_LENGTH, LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED, ZERO_REMOVAL_RATE)
		fmt.Println(p)
		fmt.Println(nerdlegame.MakeGuess("11+17=28", "19+-7=12", LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED))
		fmt.Println(nerdlegame.MakeGuess("11+17=28", "19+-7=12", LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED))
	}
}
