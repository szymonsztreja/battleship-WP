package main

import (
	"battleship-WP/apicaller"
	"battleship-WP/game"
)

func main() {
	game := game.Game{}
	game.Run()
	// apicaller.Call()
}

func dupa() {
	apicaller.Call()
	// game := game.Game{}
	// game.Run()
}
