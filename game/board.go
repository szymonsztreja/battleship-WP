package game

import (
	gui "github.com/grupawp/warships-gui/v2"
)

type WarshipBoard struct {
	x      int
	y      int
	Nick   *gui.Text
	Board  *gui.Text
	states [10][10]gui.State
	Desc   *gui.Board
}
