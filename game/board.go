package game

import (
	"fmt"
	"strconv"
	"strings"

	gui "github.com/grupawp/warships-gui/v2"
)

type WarshipBoard struct {
	x      int
	y      int
	Nick   *gui.Text
	Board  *gui.Board
	states [10][10]gui.State
	Desc   []*gui.Text
}

type coordsToCheck struct {
	x       int
	y       int
	visited bool
}

type coords struct {
	x int
	y int
}

func NewWarshipBoard(x int, y int, xDesc int, c *gui.BoardConfig) *WarshipBoard {
	wb := new(WarshipBoard)
	wb.x = x
	wb.y = y
	wb.Nick = gui.NewText(x, 30, "", nil)
	wb.Board = gui.NewBoard(x, y, c)
	setArrayValue(&wb.states, gui.Empty)
	wb.Board.SetStates(wb.states)
	// wb.Desc = gui.NewText(xDesc, 30, "", nil)
	wb.Desc = []*gui.Text{}

	return wb
}

//  TODO

func wrapText(text string, lineLength int) []string {
	words := strings.Split(text, " ")
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 > lineLength {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (wb *WarshipBoard) SetDescText(desc string) {
	lines := wrapText(desc, 44)
	wb.Desc = []*gui.Text{}

	for i, line := range lines {
		wb.Desc = append(wb.Desc, gui.NewText(wb.x, 30+i*1, line, nil))
	}
}

// // Update the description text with wrapping
// func (wb *WarshipBoard) SetDescText(desc string) {
// 	lines := wrapText(desc, 20)
// 	wb.Desc = []*gui.Text{}

// 	for i, line := range lines {
// 		wb.Desc = append(wb.Desc, gui.NewText(wb.x, 30+i*2, line, nil)) // Adjusting y-coordinate for each line
// 	}
// }

func (wb *WarshipBoard) UpdateSunk(coord string, state gui.State) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error updating sunk x:%v, y:%v", x, y)
	}

	wb.states[x][y] = state
	wb.Board.SetStates(wb.states)

}

// markForbiddenArea marks the area around a sunk ship
func (wb *WarshipBoard) UpSunk(coord string, statesToCheck []coordsToCheck) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("Error updating sunk x:%v, y:%v\n", x, y)
		return
	}

	if statesToCheck == nil {
		statesToCheck = []coordsToCheck{}
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
						sunkState := coordsToCheck{nx, ny, false}
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

// Sprawdzenie, czy koordynaty są w granicach planszy
func (wb *WarshipBoard) IsWithinBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 10 && y < 10
}

// Sprawdzenie, czy koordynaty są puste
func (wb *WarshipBoard) IsEmpty(x, y int) bool {
	return wb.IsWithinBounds(x, y) && wb.states[y][x] == gui.Empty
}

// func (wb *WarshipBoard) IsPlacementValid(coords []string) bool {
// 	for _, coord := range coords {
// 		x, y, err := stringCoordToInt(coord)
// 		if err != nil {
// 			fmt.Printf("Error placement validation %d, %d\n", x, y)
// 		}
// 		if !wb.IsEmpty(x, y) {
// 			return false
// 		}
// 		if !wb.HasAdjacentShip(x, y) {
// 			return false
// 		}
// 	}
// 	return true
// }

/*
Check adjacent

	(0,-1)

(-1,0)	x	(1,0)

	(0,1)
*/
func (wb *WarshipBoard) HasAdjacentShip(x, y int) (bool, []coords) {
	dirs := []struct{ dx, dy int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	adjacent := []coords{}

	for _, dir := range dirs {
		nx, ny := x+dir.dx, y+dir.dy
		if wb.IsWithinBounds(nx, ny) && wb.states[nx][ny] == gui.Ship {
			adjacent = append(adjacent, coords{nx, ny})
			return true, adjacent
		}
	}
	return false, adjacent
}

/*
(-1,-1)   (1,-1)

	x

(-1, 1)   (1,1)
*/
func (wb *WarshipBoard) checkDiagonally(x, y int) (bool, []coords) {
	// Define the diagonal directions
	dirs := []struct{ dx, dy int }{
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	diagonals := []coords{}

	// Check each diagonal direction
	for _, dir := range dirs {
		nx, ny := x+dir.dx, y+dir.dy
		if wb.IsWithinBounds(nx, ny) && wb.states[nx][ny] == gui.Ship {
			diagonals = append(diagonals, coords{nx, ny})
			return true, diagonals
		}
	}
	return false, diagonals
}

// Check if the ship shape is valid
func (wb *WarshipBoard) isValidShape(coord string) bool {
	x, y := GetIntCoord(coord)

	isValid := false

	hasAdjacent, adjacentCoords := wb.HasAdjacentShip(x, y)
	hasDiagonal, diagonalCoords := wb.checkDiagonally(x, y)

	// fmt.Printf("hasAdjacent: %v, adjacentCoords: %v\n", hasAdjacent, adjacentCoords)
	// fmt.Printf("hasDiagonal: %v, diagonalCoords: %v\n", hasDiagonal, diagonalCoords)
	if hasAdjacent {
		isValid = true
		return isValid
	} else if !hasAdjacent && !hasDiagonal {
		// fmt.Println("Condition: !hasAdjacent && !hasDiagonal")
		isValid = true
	} else if !hasAdjacent && hasDiagonal {
		// fmt.Println("Condition: !hasAdjacent && hasDiagonal")
		isValid = false
		// } else if hasAdjacent && !hasDiagonal {
		// 	// fmt.Println("Condition: hasAdjacent && !hasDiagonal")
		// 	isValid = true
	} else if hasAdjacent && hasDiagonal {
		fmt.Println("Condition: hasAdjacent && hasDiagonal")
		// for _, ac := range adjacentCoords {
		// Check if coords that are diagonal to our main coords
		// have any common adjacent coord
		for _, dc := range diagonalCoords {
			dcHasAdjc, dcAdjcs := wb.HasAdjacentShip(dc.x, dc.y)
			if !dcHasAdjc {
				isValid = false
			} else {
				for _, dcAdjc := range dcAdjcs {
					for _, ac := range adjacentCoords {
						if dcAdjc.x == ac.x && dcAdjc.y == ac.y {
							// fmt.Printf("Matching coordinates found: ac = %v, dc = %v\n", ac, dc)
							isValid = true
						}
					}
				}

			}
		}
		// }
	}
	// fmt.Printf("isValid: %v\n", isValid)
	return isValid
}

func (wb *WarshipBoard) getShipSize(coord string) int {
	x, y := GetIntCoord(coord)
	if !wb.IsWithinBounds(x, y) {
		return 0
	}

	// Use a stack for DFS
	stack := []coords{{x, y}}
	visited := map[coords]bool{{x, y}: true}
	size := 0

	// All direction around the ship adjacent and digonal combined
	dirs := []struct{ dx, dy int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	for len(stack) > 0 {
		// fmt.Printf("stack size %v", stack)
		c := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		size++

		for _, dir := range dirs {
			nx, ny := c.x+dir.dx, c.y+dir.dy
			neighbor := coords{nx, ny}

			if wb.IsWithinBounds(nx, ny) && wb.states[nx][ny] == gui.Ship && !visited[neighbor] {
				stack = append(stack, neighbor)
				visited[neighbor] = true
			}
		}
	}

	// fmt.Printf("Ship size at (%d, %d): %d\n", x, y, size)
	return size
}

func intCoordToString(x int, y int) string {
	column := string(rune(x + 'A'))
	row := fmt.Sprint(y + 1)
	return column + row
}

func GetIntCoord(coord string) (int, int) {
	x, y, err := stringCoordToInt(coord)
	if err != nil {
		fmt.Printf("err getting int coord: %v", err.Error())
		return 0, 0
	}
	return x, y
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
	return []gui.Drawable{wb.Nick, wb.Board}
}

func setArrayValue(arr *[10][10]gui.State, value gui.State) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			arr[i][j] = value
		}
	}
}
