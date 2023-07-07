package go_stratego

import (
	"fmt"
	"math/rand"

	wr "github.com/mroth/weightedrand"
)

type Board struct {
	board [BoardSize][BoardSize]*Unit
}

func NewEmptyBoard() *Board {
	board := [BoardSize][BoardSize]*Unit{}
	for _, pair := range [][]int{{4, 2}, {4, 3}, {5, 2}, {5, 3}, {4, 6}, {4, 7}, {5, 6}, {5, 7}} {
		board[pair[0]][pair[1]] = Water()
	}
	return &Board{
		board: board,
	}
}

func (b *Board) possibleMoves(row, col int) [][]int {
	if row < 0 || row >= BoardSize || col < 0 || col >= BoardSize {
		return make([][]int, 0)
	}
	unit := b.board[row][col]
	if unit == nil || unit.team == nil {
		return make([][]int, 0)
	}
	if unit.typ == "bomb" || unit.typ == "flag" {
		return make([][]int, 0)
	}
	moves := make([][]int, 0)
	if row+1 < BoardSize && (b.board[row+1][col] == nil || (b.board[row+1][col].team != nil && *b.board[row+1][col].team != *unit.team)) {
		moves = append(moves, []int{1, 0})
	}
	if row-1 >= 0 && (b.board[row-1][col] == nil || (b.board[row-1][col].team != nil && *b.board[row-1][col].team != *unit.team)) {
		moves = append(moves, []int{-1, 0})
	}
	if col+1 < BoardSize && (b.board[row][col+1] == nil || (b.board[row][col+1].team != nil && *b.board[row][col+1].team != *unit.team)) {
		moves = append(moves, []int{0, 1})
	}
	if col-1 >= 0 && (b.board[row][col-1] == nil || (b.board[row][col-1].team != nil && *b.board[row][col-1].team != *unit.team)) {
		moves = append(moves, []int{0, -1})
	}
	return moves
}

// retrieves the number of active (movable) units for a given team
func (b *Board) numActive(team string) int {
	count := 0
	for _, row := range b.board {
		for _, unit := range row {
			if unit != nil && *unit.team == team && unit.typ != "bomb" && unit.typ != "flag" {
				count += 1
			}
		}
	}
	return count
}

func NewRandomBoard(teams []string, random *rand.Rand) (*Board, error) {
	if len(teams) != 2 {
		return nil, fmt.Errorf("teams list must contain two teams")
	}
	teamOneUnits := map[string]int{
		"flag": 1, "bomb": 6, "spy": 1, "scout": 8, "miner": 5, "sergeant": 4, "lieutenant": 4,
		"captain": 4, "major": 3, "colonel": 2, "general": 1, "marshal": 1,
	}
	teamTwoUnits := map[string]int{
		"flag": 1, "bomb": 6, "spy": 1, "scout": 8, "miner": 5, "sergeant": 4, "lieutenant": 4,
		"captain": 4, "major": 3, "colonel": 2, "general": 1, "marshal": 1,
	}
	board := NewEmptyBoard()

	// flag with 20, 40, and 40 percent
	flagChooser, _ := wr.NewChooser(
		wr.Choice{Item: 2, Weight: 1},
		wr.Choice{Item: 1, Weight: 4},
		wr.Choice{Item: 0, Weight: 5},
	)
	place(board, flagChooser, random, true, NewUnit("flag", teams[0]))
	place(board, flagChooser, random, false, NewUnit("flag", teams[1]))
	teamOneUnits["flag"] -= 1
	teamTwoUnits["flag"] -= 1

	// bombs with 10, 20, 40, and 40 percent
	bombChooser, _ := wr.NewChooser(
		wr.Choice{Item: 3, Weight: 1},
		wr.Choice{Item: 2, Weight: 2},
		wr.Choice{Item: 1, Weight: 3},
		wr.Choice{Item: 0, Weight: 4},
	)
	for i := 0; i < 6; i++ {
		place(board, bombChooser, random, true, NewUnit("bomb", teams[0]))
		place(board, bombChooser, random, false, NewUnit("bomb", teams[1]))
		teamOneUnits["bomb"] -= 1
		teamTwoUnits["bomb"] -= 1
	}

	// miners in back three with 10, 40, 50 percent
	minerChooser, _ := wr.NewChooser(
		wr.Choice{Item: 2, Weight: 1},
		wr.Choice{Item: 1, Weight: 4},
		wr.Choice{Item: 0, Weight: 5},
	)
	for i := 0; i < 5; i++ {
		place(board, minerChooser, random, true, NewUnit("miner", teams[0]))
		place(board, minerChooser, random, false, NewUnit("miner", teams[1]))
		teamOneUnits["miner"] -= 1
		teamTwoUnits["miner"] -= 1
	}

	// scouts in front three with 10, 40, 50 percent
	scoutChooser, _ := wr.NewChooser(
		wr.Choice{Item: 3, Weight: 5},
		wr.Choice{Item: 2, Weight: 4},
		wr.Choice{Item: 1, Weight: 1},
	)
	for i := 0; i < 6; i++ {
		place(board, scoutChooser, random, true, NewUnit("scout", teams[0]))
		place(board, scoutChooser, random, false, NewUnit("scout", teams[1]))
		teamOneUnits["scout"] -= 1
		teamTwoUnits["scout"] -= 1
	}

	// place remainder randomly
	for row := 0; row < (BoardSize-2)/2; row++ {
		for col := 0; col < BoardSize; col++ {
			if board.board[row][col] == nil {
				for typ, amt := range teamOneUnits {
					if amt > 0 {
						board.board[row][col] = NewUnit(typ, teams[0])
						teamOneUnits[typ] -= 1
						break
					}
				}
			}
		}
	}
	for row := BoardSize/2 + 1; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if board.board[row][col] == nil {
				for typ, amt := range teamTwoUnits {
					if amt > 0 {
						board.board[row][col] = NewUnit(typ, teams[1])
						teamTwoUnits[typ] -= 1
						break
					}
				}
			}
		}
	}
	return board, nil
}

func getRandomNotTaken(board *Board, chooser *wr.Chooser, random *rand.Rand, isOne bool) (int, int) {
	row := chooser.PickSource(random).(int)
	col := rand.Intn(BoardSize)
	if !isOne {
		row = BoardSize - row - 1
	}
	if board.board[row][col] != nil {
		return getRandomNotTaken(board, chooser, random, isOne)
	}
	return row, col
}

func place(board *Board, chooser *wr.Chooser, random *rand.Rand, isOne bool, unit *Unit) {
	row, col := getRandomNotTaken(board, chooser, random, isOne)
	board.board[row][col] = unit
}
