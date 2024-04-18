package main

import (
	"battleship-WP/client"
	"fmt"
)

func main() {
	token := client.InitGame()
	sth, _ := client.Board(token)
	fmt.Println(sth)
}
