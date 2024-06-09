package game

import (
	"fmt"
	"strconv"

	gui "github.com/grupawp/warships-gui/v2"
)

type WarshipBoard struct {
	x      int
	y      int
	Nick   *gui.Text
	Board  *gui.Board
	states [10][10]gui.State
	Desc   *gui.Text
}

type sunkHelper struct {
	x       int
	y       int
	visited bool
}

func NewWarshipBoard(x int, y int, xDesc int, c *gui.BoardConfig) *WarshipBoard {
	wb := new(WarshipBoard)
	wb.x = x
	wb.y = y
	wb.Nick = gui.NewText(x, 30, "", nil)
	wb.Board = gui.NewBoard(x, y, c)
	setArrayValue(&wb.states, gui.Empty)
	wb.Board.SetStates(wb.states)
	wb.Desc = gui.NewText(xDesc, 30, "", nil)

	return wb
}

//  TODO

func (wb *WarshipBoard) UpdateSunk(coord string, state gui.State) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error updating sunk x:%v, y:%v", x, y)
	}

	wb.states[x][y] = state
	wb.Board.SetStates(wb.states)

}

// markForbiddenArea marks the area around a sunk ship
func (wb *WarshipBoard) UpSunk(coord string, statesToCheck []sunkHelper) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error updating sunk x:%v, y:%v\n", x, y)
		return
	}

	if statesToCheck == nil {
		statesToCheck = []sunkHelper{}
	}

	// Check if the current coordinate has already been processed
	for i := range statesToCheck {
		if statesToCheck[i].x == x && statesToCheck[i].y == y {
			if statesToCheck[i].visited {
				return
			}
			statesToCheck[i].visited = true
			break
		}
	}

	// Process the surrounding cells
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			nx, ny := x+dx, y+dy
			// Check if surrounding states are in board boundries
			if nx >= 0 && nx < 10 && ny >= 0 && ny < 10 {
				if wb.states[nx][ny] == gui.Empty {
					wb.states[nx][ny] = gui.Miss
				} else if wb.states[nx][ny] == gui.Hit {
					// Check if this state is already in statesToCheck
					alreadyInStatesToCheck := false
					for _, stateToCheck := range statesToCheck {
						if stateToCheck.x == nx && stateToCheck.y == ny {
							alreadyInStatesToCheck = true
							break
						}
					}
					if !alreadyInStatesToCheck {
						sunkState := sunkHelper{nx, ny, false}
						statesToCheck = append(statesToCheck, sunkState)
					}
				}
			}
		}
	}

	wb.Board.SetStates(wb.states)

	// Recursive call for any new coordinates added to statesToCheck
	for _, state := range statesToCheck {
		if !state.visited {
			newCoord := intCoordToString(state.x, state.y)
			wb.UpSunk(newCoord, statesToCheck)
		}
	}
}

func (wb *WarshipBoard) UpdateState(coord string, state gui.State) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error converting string to int board: x:%v, y:%v", x, y)
	}

	wb.states[x][y] = state
	wb.Board.SetStates(wb.states)
}

func (wb *WarshipBoard) Import(coords []string) {
	for _, coord := range coords {
		x, y, err := stringCoordToInt(coord)
		if err != nil {
			fmt.Printf("Error importing board: x:%v, y:%v", x, y)
		}
		wb.states[x][y] = gui.Ship
	}
}

func (wb *WarshipBoard) GetState(coord string) gui.State {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error getting board state err: %v", err.Error())
	}
	return wb.states[x][y]
}

func intCoordToString(x int, y int) string {
	column := string(rune(x + 'A'))
	row := fmt.Sprint(y + 1)
	return column + row
}

func stringCoordToInt(coord string) (int, int, error) {
	if len(coord) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinate length: %v", len(coord))
	}

	column := int(coord[0] - 'A')
	if column < 0 || column > 9 {
		return 0, 0, fmt.Errorf("invalid column in coordinate: %v", coord)
	}

	row, err := strconv.Atoi(string(coord[1:]))
	if err != nil || row < 1 || row > 10 {
		return 0, 0, fmt.Errorf("invalid row in coordinate: %v", coord[1:])
	}

	return column, row - 1, nil
}

func (wb *WarshipBoard) Drawables() []gui.Drawable {
	return []gui.Drawable{wb.Nick, wb.Board, wb.Desc}
}

func setArrayValue(arr *[10][10]gui.State, value gui.State) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			arr[i][j] = value
		}
	}
}
