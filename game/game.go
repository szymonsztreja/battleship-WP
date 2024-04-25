package game

import (
	"battleship-WP/client"
	"fmt"
	"time"

	board "github.com/grupawp/warships-lightgui/v2"
)

type Game struct {
}

func (g Game) Run() {
	gameClient := &client.GameClient{}
	gameStatus := &client.GameStatus{}
	boardStruct := client.BoardStruct{}
	var err error

	gameClient.InitGame()

	for {
		gameStatus, err = gameClient.Status()

		fmt.Printf("Game status response: %+v\n", gameStatus.GameStatus)

		if err != nil {
			fmt.Printf("error getting game status : %s\n", err)
		}

		if gameStatus.GameStatus == "game_in_progress" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Game status:", *gameStatus)

	boardStruct.Board, err = gameClient.Board()
	if err != nil {
		fmt.Printf("error getting game board: %s\n", err)
		return
	}

	// boardStruct.Board = boardArray

	fmt.Println("Game status:", boardStruct)

	board := board.New(board.NewConfig())
	err = board.Import(boardStruct.Board)
	if err != nil {
		fmt.Println(err)
	}
	board.Display()

	for {
		gameStatus, err = gameClient.Status()

		if err != nil {
			fmt.Printf("error getting game status : %s\n", err)
		}

		if gameStatus.GameStatus == "ended" {
			fmt.Print(gameStatus.LastGameStatus)
			return
		} else {
			if gameStatus.ShouldFire {

			}
		}
	}
}
