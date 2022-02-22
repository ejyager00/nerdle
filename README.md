# Nerdle

This is a recreation of the nerdle game available at [nerdlegame.com](https://nerdlegame.com/), which is itself a math-themed Wordle clone. 

You can build the project using `go build` and run the executable which should be titled `nerdle`. It will run on port 10000 by default. 

`index.html` contains a very simple frontend for it, which also displays at the root url. However, you can also query the REST API for the Go backend directly using JSON. 

The `new` endpoint will return a key for a game. It takes the parameters in the `GameStart` struct in [server.go](https://github.com/ejyager00/nerdle/blob/master/server.go). The `guess` endpoint will return several parameters, usually including a boolean for whether the guess was valid and an array indicating the successes of the guess. It takes the parameters in the `Guess` struct in [server.go](https://github.com/ejyager00/nerdle/blob/master/server.go).