package main

import (
	"math/rand"
	"strings"
)

const LETTERS string = "0123456789+-*/="

func IsValidPuzzle(puzzle string, leadingzeros, negativezeros bool) bool {
	if strings.Count(puzzle, "=") != 1 {
		return false // must have exactly one equals sign
	}
	for _, c := range puzzle {
		if !strings.ContainsRune(LETTERS, c) {
			return false // chars must be valid for puzzle
		}
	}
	if strings.Contains(puzzle, "/0") {
		return false // cannot contain zero divison
	}
	for i, c := range puzzle {
		if strings.ContainsRune("+-*/=", c) {
			if c == '-' {
				if i == 1 && strings.ContainsRune("+-*/=", rune(puzzle[0])) {
					return false // '-' cannot follow another symbol at the beginning
				} else if i > 1 && strings.ContainsRune("+-*/=", rune(puzzle[i-2])) && strings.ContainsRune("+-*/=", rune(puzzle[i-1])) {
					return false // '-' cannot end a three symbol sequence
				} else if i == len(puzzle)-1 {
					return false // cannot end with '-'
				}
			} else {
				if i == 0 {
					return false // cannot begin with non '-' symbol
				} else if strings.ContainsRune("+-*/=", rune(puzzle[i-1])) {
					return false // non '-' symbol cannot follow another symbol
				} else if i == len(puzzle)-1 {
					return false // cannot end with an operator
				}
			}
		} else {
			if c == '0' {
				if !leadingzeros && i != len(puzzle)-1 && strings.ContainsRune("0123456789", rune(puzzle[i+1])) {
					return false // no leading zeros
				}
				if !negativezeros && i > 1 && puzzle[i-1] == '-' && strings.ContainsRune("+-*/=", rune(puzzle[i-2])) {
					return false // no negative zero
				}
			}
		}
	}
	sides := strings.Split(puzzle, "=")
	left_side := sides[0]
	if left_side[0] == '-' {
		left_side = left_side[1:]
	}
	has_operator := false
	for _, c := range "+-*/" {
		if strings.ContainsRune(left_side, c) {
			has_operator = true
			break
		}
	}
	if !has_operator {
		return false // there must be an operator on the left side
	}
	has_operator = false
	right_side := sides[1]
	if right_side[0] == '-' {
		right_side = right_side[1:]
	}
	for _, c := range "+-*/" {
		if strings.ContainsRune(right_side, c) {
			has_operator = true
			break
		}
	}
	if has_operator {
		return false // there must not be an operator on the right side
	}
	result, err := IsEqual(puzzle)
	if !result || err != nil {
		return false // the equation must be true
	}
	return true
}

func RandomPuzzle(length int, leadingzeros, negativezeros bool) string {
	// create strings of random characters until one is a valid puzzle
	var puzzle string = ""
	for !IsValidPuzzle(puzzle, leadingzeros, negativezeros) {
		puzzle = ""
		for i := 0; i < length; i++ {
			puzzle += string(rune(LETTERS[rand.Intn(len(LETTERS))]))
		}
	}
	return puzzle
}

func WeightedRandomPuzzle(length int, leadingzeros, negativezeros bool, zeroremovalrate float32) string {
	// run the random puzzle function, and if a zero removal case occurs, attempt of fail and try another puzzle
	var puzzle string = ""
	new_puzzle := true
	for new_puzzle {
		new_puzzle = false
		puzzle = RandomPuzzle(length, leadingzeros, negativezeros)
		if strings.ContainsRune(puzzle, '0') {
			for i, c := range puzzle {
				if c == '0' &&
					((i == 0 || strings.ContainsRune("+-*=", rune(puzzle[i-1]))) ||
						(i != len(puzzle)-1 && i != 0 && strings.ContainsRune("+-*/", rune(puzzle[i+1])) && strings.ContainsRune("0123456789", rune(puzzle[i-1])))) {
					if rand.Float32() < zeroremovalrate {
						new_puzzle = true
						break
					}
				}
			}
		}
	}
	return puzzle
}

func MakeGuess(guess, puzzle string) []int {
	answer := make([]int, len(guess))
	for i, c := range guess {
		if c == rune(puzzle[i]) {
			answer[i] = 1 // if the guess character equals the true value
		} else if strings.ContainsRune(puzzle, c) {
			puzzle_c := strings.Count(puzzle, string(c))
			if puzzle_c >= strings.Count(guess, string(c)) { // there are more of this character in the puzzle than the guess
				answer[i] = -1 // the character is in the puzzle, but not position i
			} else if strings.Count(guess[:i], string(c)) < strings.Count(puzzle[:i], string(c)) { // there are more of this in the preceding puzzle than in the preceding guess
				answer[i] = -1 // the character is in the puzzle, but not position i
			} else if puzzle_c > strings.Count(guess[:i], string(c)) { // there are more of this character in the puzzle than the guess up to this character
				count_correct := 0
				count_remaining := strings.Count(puzzle[i:], string(c))
				for j, d := range guess[i:] {
					if d == c && rune(puzzle[j+i]) == c {
						count_correct++
					}
				}
				if count_correct == count_remaining { // all the remaining instances are "claimed" by an exact match
					answer[i] = 0 // there are not this many of this character in the puzzle
				} else {
					answer[i] = -1 // the character is in the puzzle, but not position i
				}
			} else {
				answer[i] = 0 // there are not this many of this character in the puzzle
			}
		} else {
			answer[i] = 0 // the character is not in the puzzle
		}
	}
	return answer
}
