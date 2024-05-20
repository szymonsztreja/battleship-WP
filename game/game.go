package game

import (
	"battleship-WP/client"
	"context"
	"fmt"
	"net/http"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

// var color := gui.NewColor(232, 139, 0)

// gui.Color{
// 	Red: 232
// 	Green:139
// 	Blue:0
// }

type Game struct {
	// playerStates *[10][10]gui.State
}

func (Game) Run() {
	httpClient := &client.HttpGameClient{
		Client: &http.Client{},
	}

	httpClient.InitGame()
	waitForGame(httpClient)

	desc, err := httpClient.GetPlayersDescription()
	if err != nil {
		fmt.Print(err)
	}

	playerShips := getBoardGame(httpClient)
	ui := gui.NewGUI(true)

	playerBoard := NewWarshipBoard(5, 5, nil)
	enemyBoard := NewWarshipBoard(55, 5, nil)

	playerBoard.Import(playerShips)
	drawables := append(playerBoard.Drawables(), enemyBoard.Drawables()...)

	for _, draw := range drawables {
		ui.Draw(draw)
	}

	playerBoard.Nick.SetText(desc.Nick)
	playerBoard.Desc.SetText(desc.Desc)

	enemyBoard.Nick.SetText(desc.Opponent)
	enemyBoard.Desc.SetText(desc.OppDesc)

	ui.Draw(playerBoard.Nick)
	ui.Draw(playerBoard.Desc)

	ui.Draw(enemyBoard.Nick)
	ui.Draw(enemyBoard.Desc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			status := getGameStatus(httpClient)
			turn := gui.NewText(3, 3, "Your turn!", nil)
			if status.GameStatus == "ended" {
				endText := gui.NewText(27, 33, "Game ended", nil)
				ui.Draw(endText)
				ui.Log("Game ended")
				return
			}
			if status.ShouldFire {
				handleOppShots(status.OppShots, playerBoard)
				ui.Draw(playerBoard.Board)
				ui.Draw(turn)
				handlePlayerShots(httpClient, ctx, enemyBoard, ui)
				ui.Draw(enemyBoard.Board)
			} else {
				time.Sleep(1 * time.Second)
				ui.Remove(turn)
			}
		}
	}()

	ui.Start(ctx, nil)
}

// func handlePlayerShots(httpClient *client.HttpGameClient, ctx context.Context, eb *WarshipBoard, ui *gui.GUI) {
// 	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
// 	defer cancel()

// 	coord := eb.Board.Listen(ctx)

// 	ui.Log(string(eb.GetState(coord)))
// 	fmt.Print(eb.GetState(coord))

// 	fireResponse, err := httpClient.Fire(coord)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	switch fireResponse.Result {
// 	case "hit":
// 		eb.UpdateState(coord, gui.Hit)
// 	case "miss":
// 		eb.UpdateState(coord, gui.Miss)
// 	case "sunk":
// 		eb.UpdateState(coord, gui.Hit)
// 	}

// }

func handlePlayerShots(httpClient *client.HttpGameClient, ctx context.Context, eb *WarshipBoard, ui *gui.GUI) {

	// Create a context with a timeout of 60 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var coord string
	var incorrectInput *gui.Text

	for {
		// Listen for a coordinate input from the player
		coord = eb.Board.Listen(ctx)
		if eb.GetState(coord) != gui.Empty && incorrectInput == nil {
			incorrectInput = gui.NewText(30, 35, "Please, click on an empty field!", nil)
			ui.Draw(incorrectInput)
		} else if eb.GetState(coord) == gui.Empty && incorrectInput != nil {
			ui.Remove(incorrectInput)
			break
		} else {
			break
		}
	}

	ui.Log(string(eb.GetState(coord)))

	fireResponse, err := httpClient.Fire(coord)
	if err != nil {
		fmt.Println(err)
	}

	// Update the board based on the fire response
	switch fireResponse.Result {
	case "hit":
		eb.UpdateState(coord, gui.Hit)
	case "miss":
		eb.UpdateState(coord, gui.Miss)
	case "sunk":
		eb.UpdateState(coord, gui.Hit)
	}

}

func handleOppShots(oppShots []string, pb *WarshipBoard) {
	for _, shot := range oppShots {
		if pb.GetState(shot) == gui.Ship {
			pb.UpdateState(shot, gui.Hit)
		} else {
			pb.UpdateState(shot, gui.Miss)
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

// func setPlayerBoard(coords []string) [10][10]gui.State {
// 	states := [10][10]gui.State{}
// 	for i := range states {
// 		states[i] = [10]gui.State{}
// 	}

// 	for _, coord := range coords {
// 		x, y := stringCoordToInt(coord)
// 		states[x][y] = gui.Ship
// 	}

// 	return states
// }

// Letters - rows, numbers - columns
// func stringCoordToInt(coord string) (int, int, error) {
// 	// stringCoords := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

// 	column := int(coord[0] - 'A')

// 	row, err := strconv.Atoi(coord[1:])
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, 0, err
// 	}
// 	row--

// 	return column, row, nil
// }

// A B C D E F G H I J

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
