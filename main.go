package main

import (
	"battleship-WP/apicaller"
	"battleship-WP/game"
)

func main() {
	game := game.Game{}
	game.Run()
}

func dupa() {
	apicaller.Call()
}
