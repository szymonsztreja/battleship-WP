package game

import (
	"context"
	"fmt"

	gui "github.com/grupawp/warships-gui/v2"
)

// TODO go back after setting all the ships

func PlaceShips() []string {
	ui := gui.NewGUI(true)

	txt := gui.NewText(1, 1, "Press on any coordinate to save it.", nil)
	ui.Draw(txt)
	ui.Draw(gui.NewText(1, 2, "Press Ctrl+C to exit", nil))

	coords := []string{}
	board := NewWarshipBoard(5, 20, 5, nil)

	drawables := append(board.Drawables(), board.Drawables()...)

	for _, draw := range drawables {
		ui.Draw(draw)
	}

	fmt.Println("DUppa")

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Example: Place 5 ships
	shipsToPlace := 10
	// shipsChan := make(chan string, shipsToPlace)
	// areShipsSet := false

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if shipsToPlace > 0 {
					coord := board.Board.Listen(ctx)
					ui.Log(coord)

					if coord == "" {
						incorrectInput := gui.NewText(0, 0, "Invalid coordinate!", nil)
						ui.Draw(incorrectInput)
						continue
					}

					if board.GetState(coord) != gui.Empty {
						incorrectInput := gui.NewText(0, 0, "Please, click on an empty field!", nil)
						ui.Draw(incorrectInput)
					} else {
						incorrectInput := gui.NewText(0, 0, "", nil)
						ui.Remove(incorrectInput)
						board.UpdateState(coord, gui.Ship)
						coords = append(coords, coord)
						shipsToPlace--
					}
				} else {
					return
				}
			}
		}
	}(ctx)
	ui.Start(ctx, nil)
	return coords
}
