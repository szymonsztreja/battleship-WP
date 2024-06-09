package game

import (
	"battleship-WP/client"
	"context"
	"fmt"
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
	// Coords            []string
	PlayerNick        string
	PlayerDescription string
	TargetNick        string
	Wpbot             bool
	HttpGameC         *client.HttpGameClient
}

func (g *Game) Run() {
	desc, err := g.HttpGameC.GetPlayersDescription()
	if err != nil {
		fmt.Print(err)
	}

	playerShips := getBoardGame(g.HttpGameC)
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

	// done := make(chan struct{})

	go handlePlayerShots(ctx, g.HttpGameC, enemyBoard, ui)
	go func(ctx context.Context) {
		for {
			status := getGameStatus(g.HttpGameC)
			if status.GameStatus == "ended" {
				endText.SetText(status.LastGameStatus)
				ui.Draw(playerBoard.Board)
				ui.Draw(enemyBoard.Board)
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
	}(ctx)
	ui.Start(ctx, nil)
}

func handlePlayerShots(ctx context.Context, c *client.HttpGameClient, enemyBoard *WarshipBoard, ui *gui.GUI) {
	var coord string
	incorrectInput := gui.NewText(30, 35, "", nil)
	accuracyText := gui.NewText(65, 3, "", nil)

	var shotsMissed float32 = 0
	var shotsHit float32 = 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for {
				coord = enemyBoard.Board.Listen(ctx)
				ui.Log(coord)
				if coord == "" {
					incorrectInput.SetText("Invalid coordinate!")
					ui.Draw(incorrectInput)
					continue
				}
				if enemyBoard.GetState(coord) != gui.Empty {
					incorrectInput.SetText("Please, click on an empty field!")
					ui.Draw(incorrectInput)
				} else {
					incorrectInput.SetText("")
					break
				}
			}
			// ui.Log(string(enemyBoard.GetState(coord)))

			fireResponse, err := c.Fire(coord)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch fireResponse.Result {
			case "hit":
				shotsHit++
				enemyBoard.UpdateState(coord, gui.Hit)
			case "miss":
				shotsMissed++
				enemyBoard.UpdateState(coord, gui.Miss)
			case "sunk":
				shotsHit++
				enemyBoard.UpdateState(coord, gui.Hit)
				enemyBoard.UpSunk(coord, nil)
			}

			// Set and display accuracy statistic on screen
			var accuracyString string
			if shotsMissed == 0 {
				accuracyString = "N/A"
			} else {
				shotsTaken := shotsHit + shotsMissed
				accuracy := (shotsHit / shotsTaken) * 100
				accuracyString = fmt.Sprintf("Accuracy %.2f%% (%v/%v)", accuracy, shotsHit, shotsTaken)
			}
			accuracyText.SetText(accuracyString)
			ui.Draw(accuracyText)
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

func getGameStatus(c *client.HttpGameClient) *client.GameStatus {
	var err error
	gameStatus, err := c.Status()

	if err != nil {
		fmt.Printf("error getting game status : %s\n", err)
	}
	return gameStatus
}

func getBoardGame(c *client.HttpGameClient) []string {
	ships, err := c.Board()
	if err != nil {
		fmt.Printf("error getting game board: %s\n", err)
	}
	return ships
}

func NewGame(playerNick, playerDescription, targetNick string, wpbot *bool) *Game {
	// Set default values
	// defaultPlayerCoords := []
	defaultPlayerNick := ""
	defaultPlayerDescription := ""
	defaultTargetNick := ""
	defaultWpbot := false

	// Use provided values if they are set, otherwise use defaults
	if playerNick == "" {
		playerNick = defaultPlayerNick
	}
	if playerDescription == "" {
		playerDescription = defaultPlayerDescription
	}
	if targetNick == "" {
		targetNick = defaultTargetNick
	}
	if wpbot == nil {
		wpbot = &defaultWpbot
	}

	return &Game{
		PlayerNick:        playerNick,
		PlayerDescription: playerDescription,
		TargetNick:        targetNick,
		Wpbot:             *wpbot,
	}
}
