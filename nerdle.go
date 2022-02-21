package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ejyager00/nerdle/nerdlegame"
)

const PUZZLE_LENGTH int = 8
const MAX_GUESSES int = 6
const LEADING_ZEROS_ALLOWED bool = false
const NEGATIVE_ZERO_ALLOWED bool = false
const ZERO_REMOVAL_RATE float32 = 1

type Game struct {
	puzzle        string
	length        int
	guesses       int
	maxguesses    int
	leadingzeros  bool
	negativezeros bool
}

type GameStart struct {
	Length          int
	LeadingZeros    bool
	NegativeZeros   bool
	ZeroRemovalRate float32
	MaxGuesses      int
}

type Guess struct {
	Key   int64
	Guess string
}

var games map[int64]Game
var currentKey int64 = 0

func main() {
	for i := 0; i < 20; i++ {
		p := nerdlegame.WeightedRandomPuzzle(PUZZLE_LENGTH, LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED, ZERO_REMOVAL_RATE)
		fmt.Println(p)
		fmt.Println(nerdlegame.MakeGuess("11+17=28", "19+-7=12", LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED))
		fmt.Println(nerdlegame.MakeGuess("11+17=28", "19+-7=12", LEADING_ZEROS_ALLOWED, NEGATIVE_ZERO_ALLOWED))
	}
}

func handleRequests() {
	http.HandleFunc("/", home)
	http.HandleFunc("/new", newPuzzle)
	http.HandleFunc("/guess", guessHandle)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Hello, world!")
	fmt.Println("Endpoint Hit: homePage")
}

func newPuzzle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	var req GameStart
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		puzzle := nerdlegame.WeightedRandomPuzzle(req.Length, req.LeadingZeros, req.NegativeZeros, req.ZeroRemovalRate)
		game := Game{puzzle: puzzle, length: req.Length, guesses: 0, maxguesses: req.MaxGuesses, leadingzeros: req.LeadingZeros, negativezeros: req.NegativeZeros}
		key := currentKey
		currentKey++
		games[key] = game
		fmt.Printf("{\"key\":%d}\n", key)
		fmt.Fprintf(w, "{\"key\":%d}", key)
	}
}

func guessHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	var req Guess
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	game := games[req.Key]
	if !nerdlegame.IsValidPuzzle(req.Guess, game.leadingzeros, game.negativezeros) || len(req.Guess) != game.length {
		fmt.Fprintf(w, "{\"validguess\":false}")
	} else {
		game.guesses++
		comparison := nerdlegame.MakeGuess(req.Guess, game.puzzle, game.leadingzeros, game.negativezeros)
		won := true
		for _, c := range comparison {
			if c != 1 {
				won = false
			}
		}
		var comparison_string string = "[" + fmt.Sprint(comparison[0])
		for _, c := range comparison[1:] {
			comparison_string += ", " + fmt.Sprint(c)
		}
		comparison_string += "]"
		if won {
			fmt.Fprintf(w, "{\"comparison\"%s:,\"validguess\":true,\"won\":true,\"loss\":false,\"guesses\":%d}", comparison_string, game.guesses)
		} else if game.guesses == game.maxguesses {
			fmt.Fprintf(w, "{\"comparison\":%s,\"validguess\":true,\"won\":false,\"loss\":true,\"guesses\":%d,\"solution\":%#v}", comparison_string, game.guesses, game.puzzle)
		} else {
			fmt.Fprintf(w, "{\"comparison\":%s,\"validguess\":true,\"won\":false,\"loss\":false,\"guesses\":%d}}", comparison_string, game.guesses)
		}
		games[req.Key] = game
	}
}
