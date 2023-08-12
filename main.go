package main

import "github.com/EdsonHTJ/stockfish-api/router"

func main() {
	router.New().Run(":8080")
}
