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
	incorrectInput := gui.NewText(30, 3, "", nil)

	coords := []string{}
	board := NewWarshipBoard(5, 20, 5, nil)

	drawables := append(board.Drawables(), board.Drawables()...)

	for _, draw := range drawables {
		ui.Draw(draw)
	}

	fmt.Println("DUppa")

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	shipsToPlace := 10
	areShipsSet := false

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
						incorrectInput.SetText("Invalid coordinate!")
						ui.Draw(incorrectInput)
						continue
					}
					state := board.GetState(coord)
					// If state is other than gui.Miss
					if state != gui.Empty && state != gui.Ship {
						incorrectInput.SetText("Please, click on an empty field!")
						ui.Draw(incorrectInput)
						// If ship is already clicked, click on it to unset it
					} else if state == gui.Ship {
						board.UpdateState(coord, gui.Empty)
						// Set empty state to ship
					} else {
						x, y := GetIntCoord(coord)
						if board.checkDiagonally(x, y) && !board.HasAdjacentShip(x, y) {
							incorrectInput.SetText("Ships cannot be diagonally!")
							ui.Log(fmt.Sprintf("Ships cannot be diagonally: coord:%v!", coord))
							break
						}
						incorrectInput := gui.NewText(0, 0, "", nil)
						ui.Remove(incorrectInput)
						board.UpdateState(coord, gui.Ship)
						coords = append(coords, coord)
						shipsToPlace--
					}
				} else {
					areShipsSet = true
				}
			}
			if areShipsSet {
				ui.Log("Ships set!")
				return
			}
		}
	}(ctx)
	ui.Start(ctx, nil)
	return coords
}
