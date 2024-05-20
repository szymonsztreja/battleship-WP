package game

import (
	"battleship-WP/client"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

type Game struct {
	// playerStates *[10][10]gui.State
}

func (Game) Run() {
	httpClient := &client.HttpGameClient{
		Client: &http.Client{},
	}
	// gameStatus := &client.GameStatus{}
	// playerStates := &[10][10]gui.State{}
	// var err error

	httpClient.InitGame()
	waitForGame(httpClient)

	desc, err := httpClient.GetPlayersDescription()
	if err != nil {
		fmt.Print(err)
	}

	playerShips := getBoardGame(httpClient)

	ui := gui.NewGUI(true)
	playerBoard := gui.NewBoard(5, 5, nil)
	enemyBoard := gui.NewBoard(55, 5, nil)

	playerStates := setPlayerBoard(playerShips)
	enemyStates := [10][10]gui.State{}
	for i := range enemyStates {
		enemyStates[i] = [10]gui.State{}
	}

	playerBoard.SetStates(playerStates)
	enemyBoard.SetStates(enemyStates)

	ui.Draw(playerBoard)
	ui.Draw(enemyBoard)

	playerNick := gui.NewText(1, 1, desc.Nick, nil)
	enemyNick := gui.NewText(55, 1, desc.Opponent, nil)
	pDesc := gui.NewText(1, 30, desc.Desc, nil)
	eDesc := gui.NewText(55, 30, desc.OppDesc, nil)
	ui.Draw(playerNick)
	ui.Draw(enemyNick)

	ui.Draw(pDesc)
	ui.Draw(eDesc)

	ctx := context.Background()

	go func() {
		for {
			status := getGameStatus(httpClient)
			turn := gui.NewText(3, 3, "Your turn!", nil)
			fmt.Print("duppppa")
			if status.GameStatus == "ended" {
				fmt.Print("Game ended")
				ui.Log("Game ended")
				return
			}
			if status.ShouldFire {
				handleOppShots(status.OppShots, &playerStates)
				// updateDisplay()
				playerBoard.SetStates(playerStates)
				ui.Draw(playerBoard)
				ui.Draw(turn)
				fmt.Print("I was here")
				handlePlayerShots(httpClient, ctx, enemyBoard, &enemyStates)
				coord := enemyBoard.Listen(ctx)

				fireResponse, err := httpClient.Fire(coord)
				if err != nil {
					ui.Log(string(err.Error()))
				}

				x, y := stringCoordToInt(coord)

				switch fireResponse.Result {
				case "hit":
					enemyStates[x][y] = gui.Hit
				case "miss":
					enemyStates[x][y] = gui.Miss
				case "sunk":
					enemyStates[x][y] = gui.Hit
				}
				enemyBoard.SetStates(enemyStates)
				ui.Draw(enemyBoard)
			} else {
				time.Sleep(1 * time.Second)
				ui.Remove(turn)
			}
		}
	}()
	ui.Start(ctx, nil)
}

// func setupBoards(httpClient *client.HttpGameClient) {
// 	playerShips := getBoardGame(httpClient)

// 	ui := gui.NewGUI(true)
// 	playerBoard := gui.NewBoard(5, 5, nil)
// 	enemyBoard := gui.NewBoard(55, 5, nil)

// 	playerStates := setPlayerBoard(playerShips)
// 	enemyStates := [10][10]gui.State{}
// 	for i := range enemyStates {
// 		enemyStates[i] = [10]gui.State{}
// 	}

// 	playerBoard.SetStates(playerStates)
// 	enemyBoard.SetStates(enemyStates)

// 	ui.Draw(playerBoard)
// 	ui.Draw(enemyBoard)

// 	ctx := context.Background()
// 	ui.Start(ctx, nil)
// }

func handlePlayerShots(httpClient *client.HttpGameClient, ctx context.Context, enemyBoard *gui.Board, enemyStates *[10][10]gui.State) {

	// Start the timer
	timer := time.NewTimer(60 * time.Second)

	// Channel to signal when time is up
	timerCh := make(chan bool)

	// Goroutine to wait for the timer to expire
	go func() {
		<-timer.C
		fmt.Println("Time's up!")
		timerCh <- true
	}()

	select {
	case <-timerCh:
		fmt.Println("Shooting time expired")
		// You can handle the case when the time expires here
		return
	case <-ctx.Done():
		fmt.Println("Context canceled")
		// Handle cancellation of context
		return
	case <-ctx.Done():
		fmt.Println("Context canceled")
		// Handle cancellation of context
		// return
	}

	coord := enemyBoard.Listen(ctx)
	fireResponse, err := httpClient.Fire(coord)
	if err != nil {
		fmt.Println(err)
	}

	x, y := stringCoordToInt(coord)

	switch fireResponse.Result {
	case "hit":
		enemyStates[x][y] = gui.Hit
	case "miss":
		enemyStates[x][y] = gui.Miss
	case "sunk":
		enemyStates[x][y] = gui.Hit
	}
}

func handleOppShots(oppShots []string, pStates *[10][10]gui.State) {
	for _, shot := range oppShots {
		x, y := stringCoordToInt(shot)
		if pStates[x][y] == gui.Ship {
			pStates[x][y] = gui.Hit
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

func setPlayerBoard(coords []string) [10][10]gui.State {
	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}

	for _, coord := range coords {
		x, y := stringCoordToInt(coord)
		states[x][y] = gui.Ship
	}

	return states
}

// Letters - rows, numbers - columns
func stringCoordToInt(coord string) (int, int) {
	// stringCoords := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	column := int(coord[0] - 'A')

	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		fmt.Println(err)
	}
	row--

	return column, row
}

// A B C D E F G H I J
