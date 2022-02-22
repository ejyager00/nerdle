package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

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
	games = make(map[int64]Game)
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/", home)
	http.HandleFunc("/new", newPuzzle)
	http.HandleFunc("/guess", guessHandle)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	index, _ := os.ReadFile("index.html")
	fmt.Fprint(w, string(index))
}

func newPuzzle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var req GameStart
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		puzzle := WeightedRandomPuzzle(req.Length, req.LeadingZeros, req.NegativeZeros, req.ZeroRemovalRate)
		game := Game{puzzle: puzzle, length: req.Length, guesses: 0, maxguesses: req.MaxGuesses, leadingzeros: req.LeadingZeros, negativezeros: req.NegativeZeros}
		key := currentKey
		currentKey++
		games[key] = game
		fmt.Fprintf(w, "{\"key\":%d}", key)
	}
}

func guessHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var req Guess
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	game := games[req.Key]
	if !IsValidPuzzle(req.Guess, game.leadingzeros, game.negativezeros) || len(req.Guess) != game.length {
		fmt.Fprintf(w, "{\"validguess\":false}")
	} else {
		game.guesses++
		comparison := MakeGuess(req.Guess, game.puzzle, game.leadingzeros, game.negativezeros)
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
			fmt.Fprintf(w, "{\"comparison\":%s,\"validguess\":true,\"won\":true,\"loss\":false,\"guesses\":%d}", comparison_string, game.guesses)
		} else if game.guesses == game.maxguesses {
			fmt.Fprintf(w, "{\"comparison\":%s,\"validguess\":true,\"won\":false,\"loss\":true,\"guesses\":%d,\"solution\":%#v}", comparison_string, game.guesses, game.puzzle)
		} else {
			fmt.Fprintf(w, "{\"comparison\":%s,\"validguess\":true,\"won\":false,\"loss\":false,\"guesses\":%d}", comparison_string, game.guesses)
		}
		games[req.Key] = game
	}
}
