package game

import (
	"context"
	"fmt"

	gui "github.com/grupawp/warships-gui/v2"
)

func PlaceShips() []string {
	ui := gui.NewGUI(true)

	ui.Draw(gui.NewText(10, 1, "Press on coordinates to set your board.", nil))
	ui.Draw(gui.NewText(10, 2, "However you have to set your ship form biggest to smallest", nil))
	ui.Draw(gui.NewText(10, 3, "Press Ctrl+C to exit", nil))
	incorrectInput := gui.NewText(30, 5, "", nil)

	coords := []string{}
	board := NewWarshipBoard(20, 7, 5, nil)

	drawables := append(board.Drawables(), board.Drawables()...)
	drawables = append(drawables, incorrectInput)

	for _, draw := range drawables {
		ui.Draw(draw)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Fleet
	fleet := map[int]int{
		4: 1, // One 4-mast ship
		3: 2, // Two 3-mast ships
		2: 3, // Three 2-mast ships
		1: 4, // Four 1-mast ships
	}

	// To keep track of placed ships
	placedShips := map[int]int{
		4: 0,
		3: 0,
		2: 0,
		1: 0,
	}

	shipsToPlace := 20
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
						// } else if state == gui.Ship {
						// 	shipSize := board.getShipSize(coord)
						// 	placedShips[shipSize]--
						// 	board.UpdateState(coord, gui.Empty)
						// Set empty state to ship
					} else {
						// x, y := GetIntCoord(coord)
						if !board.isValidShape(coord) {
							incorrectInput.SetText("Ships cannot be diagonally!")
							ui.Log(fmt.Sprintf("Ships cannot be diagonally: coord:%v!", coord))
							break
						}
						shipSize := board.getShipSize(coord)
						// Check if the ship size exists in the fleet map
						if _, ok := fleet[shipSize]; !ok {
							incorrectInput.SetText(fmt.Sprintf("Invalid ship size: %d-mast ship!", shipSize))
							ui.Log(fmt.Sprintf("Invalid ship size: %d-mast ship: coord:%v!", shipSize, coord))
							ui.Draw(incorrectInput)
							break
						}
						// Check if the all ships of particular size are set
						if placedShips[shipSize] >= fleet[shipSize] {
							incorrectInput.SetText(fmt.Sprintf("No more %d-mast ships allowed!", shipSize))
							ui.Log(fmt.Sprintf("No more %d-mast ships allowed: coord:%v!", shipSize, coord))
							ui.Draw(incorrectInput)
							break
						} else {
							incorrectInput := gui.NewText(0, 0, "", nil)
							ui.Remove(incorrectInput)
							board.UpdateState(coord, gui.Ship)
							coords = append(coords, coord)
							if shipSize > 0 {
								placedShips[shipSize-1]--
								placedShips[shipSize]++
							}
							shipsToPlace--
							ui.Log(fmt.Sprintf("Updating ship of size %v, ships_amount %v", shipSize, placedShips[shipSize]))
						}

						// shipsToPlace--
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
