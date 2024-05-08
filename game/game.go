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

func (Game) Run() {
	httpClient := &client.HttpGameClient{
		Client: &http.Client{},
	}
	gameStatus := &client.GameStatus{}
	var err error

	httpClient.InitGame()
	waitForGame(httpClient)
	yourShips := getBoardGame(httpClient)

	b := board.New(board.NewConfig())
	err = b.Import(yourShips)
	if err != nil {
		fmt.Println(err)
	}
	for {
		status := getGameStatus(httpClient)

		if status.GameStatus == "ended" {
			fmt.Print("Game ended")
			return
		}
		if status.ShouldFire {
			handleOppShots(gameStatus, b)
			b.Display()
			fmt.Println("Your turn!")
			handleYourShots(httpClient, b)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func handleYourShots(httpClient *client.HttpGameClient, b *board.Board) {
	var err error
	var prompt string
	yourShot, ok := board.ReadLineWithTimer(prompt, 60*time.Second)
	if !ok {
		fmt.Printf("There was a problem with reading line %v", ok)
	}

	fireResponse, err := httpClient.Fire(yourShot)
	if err != nil {
		fmt.Println(err)
	}
	SetRightBoard(fireResponse.Result, yourShot, b)
}

func handleOppShots(gameStatus *client.GameStatus, b *board.Board) {
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
}

func getGameStatus(httpClient *client.HttpGameClient) *client.GameStatus {
	var err error
	gameStatus, err := httpClient.Status()

	if err != nil {
		fmt.Printf("error getting game status : %s\n", err)
	}
	return gameStatus
}

func getBoardGame(httpClient *client.HttpGameClient) []string {
	ships, err := httpClient.Board()
	if err != nil {
		fmt.Printf("error getting game board: %s\n", err)
	}
	return ships
}

func waitForGame(httpClient *client.HttpGameClient) {
	for {
		status := getGameStatus(httpClient)

		if status.GameStatus == "game_in_progress" {
			break
		}

		time.Sleep(1 * time.Second)
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
