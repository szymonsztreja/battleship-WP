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
	fireResponse := &client.FireResponse{}

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

	boardStruct.Board, err = gameClient.Board()
	if err != nil {
		fmt.Printf("error getting game board: %s\n", err)
		return
	}

	fmt.Println("Game status:", boardStruct)

	b := board.New(board.NewConfig())
	err = b.Import(boardStruct.Board)
	if err != nil {
		fmt.Println(err)
	}
	b.Display()

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
				for _, shot := range gameStatus.OppShots {

					state, err := b.HitOrMiss(board.Left, shot)
					if err != nil {
						fmt.Println(err)
					}

					err = b.Set(board.Left, shot, state)
					if err != nil {
						fmt.Println(err)
					}
				}
				b.Display()
				fmt.Println("Your turn!")

				var prompt string
				yourShot, ok := board.ReadLineWithTimer(prompt, 60*time.Second)
				if ok {
				}

				fireResponse, err = gameClient.Fire(yourShot)
				if err != nil {
					fmt.Println(err)
				}

				SetRightBoard(fireResponse.Result, yourShot, b)

			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func SetRightBoard(s string, yourShot string, b *board.Board) {
	var err error
	switch s {
	case "hit":
		err = b.Set(board.Right, yourShot, board.Hit)
		if err != nil {
			fmt.Println(err)
		}
	case "miss":
		err = b.Set(board.Right, yourShot, board.Miss)
		if err != nil {
			fmt.Println(err)
		}
	case "sunk":
		err = b.Set(board.Right, yourShot, board.Hit)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// func GetState(s string){
// 	switch s {
// 	case "hit":
// 		state := board.Hit
// 	case "miss":
// 		state := board.Miss
// 	case "sunk":
// 		state := board.Hit
// 	}

// 	return state
// }
