package main

import (
	"fmt"
	"math"
)

func Generate(puzzle_length int) {
	all_puzzles := make([]string, 0, 8192)
	num_puzzles := 0
	for i := 0; i < int(math.Pow(float64(len(LETTERS)), float64(puzzle_length))); i++ {
		puzzle := ""
		for j := 0; j < puzzle_length; j++ {
			puzzle += string(rune(LETTERS[int(i/int(math.Pow(float64(len(LETTERS)), float64(j))))%len(LETTERS)]))
		}
		if isValidSolutionPuzzle(puzzle) {
			num_puzzles++
			all_puzzles = append(all_puzzles, puzzle)
			if num_puzzles%1000 == 0 {
				fmt.Println(puzzle, num_puzzles)
			}
		}
	}
	fmt.Println(all_puzzles)

}

func isValidSolutionPuzzle(puzzle string) bool {
	equal, err := IsEqual(puzzle)
	if !equal || err != nil {
		return false
	}
	lead_zero := LeadingZeros(puzzle)
	if lead_zero {
		return false
	}
	neg_zero := NegativeZeros(puzzle)
	if neg_zero {
		return false
	}
	ops := ContainsOperator(puzzle)
	if !ops[0] || ops[1] {
		return false
	}
	zero_mult, _ := MultiplicationByZero(puzzle)
	if zero_mult {
		return false
	}
	zero_div, _ := DivisionOfZero(puzzle)
	return !zero_div
}
