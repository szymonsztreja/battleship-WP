package game

import (
	"battleship-WP/client"
	"fmt"
	"net/http"
	"time"

	board "github.com/grupawp/warships-lightgui/v2"
)

type Game struct {
}

func (g Game) Run() {
	httpClient := &client.HttpGameClient{
		Client: &http.Client{},
	}
	gameStatus := &client.GameStatus{}
	boardStruct := client.BoardStruct{}
	fireResponse := &client.FireResponse{}

	var err error

	dupa := httpClient.InitGame()
	fmt.Println(dupa)
	for {
		gameStatus, err = httpClient.Status()

		fmt.Printf("Game status response: %+v\n", gameStatus.GameStatus)

		if err != nil {
			fmt.Printf("error getting game status : %s\n", err)
		}

		if gameStatus.GameStatus == "game_in_progress" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	boardStruct.Board, err = httpClient.Board()
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
		gameStatus, err = httpClient.Status()

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

				fireResponse, err = httpClient.Fire(yourShot)
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
