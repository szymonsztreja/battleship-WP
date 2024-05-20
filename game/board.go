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

func NewWarshipBoard(x int, y int, c *gui.BoardConfig) *WarshipBoard {
	wb := new(WarshipBoard)
	wb.x = x
	wb.y = y
	wb.Nick = gui.NewText(x, 30, "", nil)
	wb.Board = gui.NewBoard(x, y, c)
	setArrayValue(&wb.states, gui.Empty)
	wb.Board.SetStates(wb.states)
	wb.Desc = gui.NewText(1, 30, "", nil)

	return wb
}

func (wb *WarshipBoard) UpdateState(coord string, state gui.State) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error converting string to int board: %v", wb.Nick)
	}

	wb.states[x][y] = state
	wb.Board.SetStates(wb.states)
}

func (wb *WarshipBoard) Import(coords []string) {
	for _, coord := range coords {
		x, y, err := stringCoordToInt(coord)
		if err != nil {
			fmt.Printf("Error converting string to int board: %v", wb.Nick)
		}
		wb.states[x][y] = gui.Ship
	}
}

func (wb *WarshipBoard) GetState(coord string) gui.State {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error converting string to int board: %v", wb.Nick)
	}
	return wb.states[x][y]
}

// func ConverCoordToInt(coord string) (int, int) {
// 	x, y, err := stringCoordToInt(coord)
// 	if err != nil {
// 		fmt.Printf("Error converting string to int board: ")
// 	}
// 	return x, y
// }

func stringCoordToInt(coord string) (int, int, error) {
	// stringCoords := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	column := int(coord[0] - 'A')

	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	row--

	return column, row, nil
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
