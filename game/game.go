package game

import (
	"battleship-WP/client"
	"context"
	"fmt"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

var defaultShips = []string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"}

// var color := gui.NewColor(232, 139, 0)

// gui.Color{
// 	Red: 232
// 	Green:139
// 	Blue:0
// }

type Game struct {
	Coords            []string
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

	// playerBoard.Nick.SetText(desc.Nick)
	// playerBoard.SetDescText(desc.Desc)

	// enemyBoard.Nick.SetText(desc.Opponent)
	// enemyBoard.SetDescText(desc.OppDesc)

	// ui.Draw(playerBoard.Nick)
	// ui.Draw(playerBoard.Desc)

	// ui.Draw(enemyBoard.Nick)
	// ui.Draw(enemyBoard.Desc)

	drawNicksAndDesc(ui, desc, playerBoard, enemyBoard)

	ui.Log(g.HttpGameC.XAuthToken)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	turn := gui.NewText(47, 3, "", nil)
	timer := gui.NewText(51, 1, "", nil)
	endText := gui.NewText(51, 33, "Game ended", nil)

	// Send request to an api and cancel the game
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	// <-ctx.Done()
	// 	g.HttpGameC.AbandonGame()
	// 	cancel()
	// }()
	// go quitGame(g, ui, cancel)
	go handlePlayerShots(ctx, g.HttpGameC, enemyBoard, ui)
	go handleOppShots(ctx, g.HttpGameC, playerBoard, ui)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				status := getGameStatus(g.HttpGameC)
				if status.GameStatus == "ended" {
					endText.SetText(status.LastGameStatus)
					ui.Draw(playerBoard.Board)
					ui.Draw(enemyBoard.Board)
					ui.Draw(endText)
					ui.Log("Game ended")
					time.Sleep(2 * time.Second)
					return
				}
				if !status.ShouldFire {
					turn.SetText("Opponent turn!")
					time.Sleep(500 * time.Millisecond)
					timer.SetText("-")
					ui.Draw(timer)
				} else if status.ShouldFire {
					turn.SetText("Your turn!")
					timer.SetText(fmt.Sprint(status.Timer))
					ui.Draw(timer)
					ui.Draw(turn)
				}
			}

		}
	}(ctx)
	ui.Start(ctx, nil)

	s := getGameStatus(g.HttpGameC)
	// Send game quiting signal to an api
	if s.GameStatus == "game_in_progress" {
		g.HttpGameC.AbandonGame()
	}
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

func handleOppShots(ctx context.Context, c *client.HttpGameClient, pb *WarshipBoard, ui *gui.GUI) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// oppShotsCount = 0
			status, err := c.Status()
			if err != nil {
				ui.Log("Error getting opp shots")
			}
			oppShots := status.OppShots
			for _, shot := range oppShots {
				if pb.GetState(shot) == gui.Ship {
					pb.UpdateState(shot, gui.Hit)
					ui.Draw(pb.Board)
				} else if pb.GetState(shot) == gui.Empty {
					pb.UpdateState(shot, gui.Miss)
					ui.Draw(pb.Board)
				}
				// ui.Draw(pb.Board)
			}
		}
	}

}

func drawNicksAndDesc(ui *gui.GUI, desc *client.PlayersDescription, playerBoard *WarshipBoard, enemyBoard *WarshipBoard) {
	playerBoard.Nick.SetText(desc.Nick)
	playerBoard.SetDescText(desc.Desc)

	enemyBoard.Nick.SetText(desc.Opponent)
	enemyBoard.SetDescText(desc.OppDesc)

	ui.Draw(playerBoard.Nick)
	// ui.Draw(playerBoard.Desc)

	for _, pd := range playerBoard.Desc {
		ui.Draw(pd)
	}

	ui.Draw(enemyBoard.Nick)
	// ui.Draw(enemyBoard.Desc)

	for _, pd := range enemyBoard.Desc {
		ui.Draw(pd)
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
