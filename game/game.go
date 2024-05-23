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
	PlayerNick        string
	PlayerDescription string
}

func (g *Game) Run() {
	httpClient := &client.HttpGameClient{
		Client: &http.Client{},
	}

	gameData := client.GameData{
		Coords:     []string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"},
		Desc:       g.PlayerDescription,
		Nick:       g.PlayerNick,
		TargetNick: "",
		Wpbot:      true,
	}

	httpClient.InitGame(gameData)
	waitForGame(httpClient)

	desc, err := httpClient.GetPlayersDescription()
	if err != nil {
		fmt.Print(err)
	}

	playerShips := getBoardGame(httpClient)
	ui := gui.NewGUI(true)

	playerBoard := NewWarshipBoard(5, 5, 5, nil)
	enemyBoard := NewWarshipBoard(55, 5, 55, nil)

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

	turn := gui.NewText(47, 3, "", nil)
	timer := gui.NewText(51, 1, "", nil)
	endText := gui.NewText(51, 33, "Game ended", nil)

	go func() {
		for {
			status := getGameStatus(httpClient)
			if status.GameStatus == "ended" {
				endText.SetText("Gamee ended")
				ui.Draw(endText)
				ui.Log("Game ended")
				return
			}
			if !status.ShouldFire {
				turn.SetText("Opponent turn!")
				time.Sleep(1 * time.Second)
				timer.SetText("-")
				ui.Draw(timer)
			} else if status.ShouldFire {
				turn.SetText("Your turn!")
				timer.SetText(fmt.Sprint(status.Timer))
				ui.Draw(timer)
				ui.Draw(turn)
				handleOppShots(status.OppShots, playerBoard, ui)
			}
		}
	}()

	go handlePlayerShots(ctx, httpClient, enemyBoard, ui)

	ui.Start(ctx, nil)
}

func handlePlayerShots(ctx context.Context, httpClient *client.HttpGameClient, enemyBoard *WarshipBoard, ui *gui.GUI) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var coord string
	incorrectInput := gui.NewText(30, 35, "", nil)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for {
				coord = enemyBoard.Board.Listen(ctx)
				if enemyBoard.GetState(coord) != gui.Empty {
					incorrectInput.SetText("Please, click on an empty field!")
					ui.Draw(incorrectInput)
				} else if enemyBoard.GetState(coord) == gui.Empty {
					incorrectInput.SetText("")
					break
				}
			}
			ui.Log(string(enemyBoard.GetState(coord)))

			fireResponse, err := httpClient.Fire(coord)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch fireResponse.Result {
			case "hit":
				enemyBoard.UpdateState(coord, gui.Hit)
			case "miss":
				enemyBoard.UpdateState(coord, gui.Miss)
			case "sunk":
				enemyBoard.UpdateState(coord, gui.Hit)
			}
		}
	}
}

func handleOppShots(oppShots []string, pb *WarshipBoard, ui *gui.GUI) {
	for _, shot := range oppShots {
		if pb.GetState(shot) == gui.Ship {
			pb.UpdateState(shot, gui.Hit)
		} else {
			pb.UpdateState(shot, gui.Miss)
		}
		ui.Draw(pb.Board)
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
