package go_stratego

import (
	"fmt"
	"math"
	"math/rand"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

type state struct {
	turn        string
	teams       []string
	winners     []string
	ready       map[string]bool
	board       *Board
	battle      *Battle
	justBattled bool
	started     bool
}

func newState(teams []string, variant string, random *rand.Rand) (*state, error) {
	if random == nil {
		return nil, fmt.Errorf("random seed is null")
	}
	board, err := NewRandomBoard(teams, variant, random)
	if err != nil {
		return nil, err
	}
	ready := make(map[string]bool)
	for _, team := range teams {
		ready[team] = false
	}
	return &state{
		board:   board,
		teams:   teams,
		turn:    teams[0],
		winners: make([]string, 0),
		ready:   ready,
	}, nil
}

func (s *state) ToggleReady(team string) error {
	if s.started {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot toggle ready when game has already started"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	s.ready[team] = !s.ready[team]
	if s.playersReady() {
		s.started = true
	}
	return nil
}

func (s *state) SwitchUnits(team string, unitRow, unitCol, switchRow, switchCol int) error {
	boardSize := len(s.board.board)
	if s.started {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units when game has already started"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if s.ready[team] {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units when you are ready"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if unitRow >= boardSize || unitRow < 0 || unitCol >= boardSize || unitCol < 0 ||
		switchRow >= boardSize || switchRow < 0 || switchCol >= boardSize || switchCol < 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("index out of bounds"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	unit := s.board.board[unitRow][unitCol]
	if unit == nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("no unit at %d,%d", unitRow, unitCol),
			Status: bgerr.StatusInvalidAction,
		}
	}

	var minRow, maxRow int
	if *unit.Team == s.teams[0] {
		minRow = 0
		maxRow = (boardSize / 2) - 2
	} else {
		minRow = (boardSize / 2) + 1
		maxRow = boardSize - 1
	}

	if switchRow < minRow || switchRow > maxRow {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch unit outside of your side"),
			Status: bgerr.StatusInvalidAction,
		}
	}

	swtch := s.board.board[switchRow][switchCol]
	if unit.Team == nil || (swtch != nil && swtch.Team == nil) {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units that have no team"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if swtch != nil && *unit.Team != *swtch.Team {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units that are not on the same team"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if *unit.Team != team {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch the opposing team's units"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	s.board.board[unitRow][unitCol] = swtch
	s.board.board[switchRow][switchCol] = unit
	return nil
}

func (s *state) MoveUnit(team string, unitRow, unitCol, moveRow, moveCol int) error {
	boardSize := len(s.board.board)
	if !s.playersReady() {
		return &bgerr.Error{
			Err:    fmt.Errorf("both players are not ready"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot play on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if unitRow >= boardSize || unitRow < 0 || unitCol >= boardSize || unitCol < 0 ||
		moveRow >= boardSize || moveRow < 0 || moveCol >= boardSize || moveCol < 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("index out of bounds"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	unit := s.board.board[unitRow][unitCol]
	if unit == nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("unit does not exist at %d, %d", unitRow, unitCol),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if *unit.Team != team {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot move a unit not part of your team"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if unit.Type == bomb || unit.Type == flag {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot move bombs or flags"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if (unit.Type != scout && math.Abs(float64(moveRow)-float64(unitRow))+math.Abs(float64(moveCol)-float64(unitCol)) > 1.0) ||
		(unit.Type == scout && math.Abs(float64(moveRow)-float64(unitRow)) > 1.0 && math.Abs(float64(moveCol)-float64(unitCol)) > 1.0) {
		return &bgerr.Error{
			Err:    fmt.Errorf("unit cannot move diagonally"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if (math.Abs(float64(moveRow)-float64(unitRow)) > 1.0 && unit.Type != scout) ||
		(math.Abs(float64(moveCol)-float64(unitCol)) > 1.0 && unit.Type != scout) {
		return &bgerr.Error{
			Err:    fmt.Errorf("unit cannot move more than one space unless they are a scout"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if unit.Type == scout && !scoutCanMove(s.board, unitRow, unitCol, moveRow, moveCol, *unit.Team) {
		return &bgerr.Error{
			Err:    fmt.Errorf("scout cannot move through water or other units"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	attackedUnit := s.board.board[moveRow][moveCol]
	if attackedUnit != nil {
		if attackedUnit.Type == water {
			return &bgerr.Error{
				Err:    fmt.Errorf("cannot move onto water"),
				Status: bgerr.StatusInvalidAction,
			}
		}
		winningUnit, err := unit.Attack(attackedUnit)
		if err != nil {
			return err
		}
		s.board.board[unitRow][unitCol] = nil
		s.board.board[moveRow][moveCol] = winningUnit
		winner := ""
		if winningUnit != nil {
			winner = *winningUnit.Team
		}
		// check for game over
		if attackedUnit.Type == flag {
			s.winners = []string{team} // attacked flag so game is over
		} else if s.board.numActive(*attackedUnit.Team) == 0 && s.board.numActive(team) == 0 {
			s.winners = []string{""} // both teams ran out of movable units
		} else if s.board.numActive(*attackedUnit.Team) == 0 {
			s.winners = []string{team} // one team ran out of movable units
		} else if s.board.numActive(team) == 0 {
			s.winners = []string{*attackedUnit.Team} // the other team ran out of movable units
		}
		s.nextTurn()
		s.justBattled = true
		s.battle = &Battle{
			MoveUnitActionDetails: MoveUnitActionDetails{
				unitRow, unitCol, moveRow, moveCol,
			},
			AttackingUnit: *unit,
			AttackedUnit:  *attackedUnit,
			WinningTeam:   winner,
		}
		return nil
	} else {
		s.board.board[unitRow][unitCol] = nil
		s.board.board[moveRow][moveCol] = unit
		s.nextTurn()
		s.started = true
		s.justBattled = false
		return nil
	}
}

func scoutCanMove(board *Board, scoutRow, scoutCol, moveRow, moveCol int, scoutTeam string) bool {
	rowDirection := 0
	if moveRow-scoutRow > 0 {
		rowDirection = 1
	} else if moveRow-scoutRow < 0 {
		rowDirection = -1
	}
	if rowDirection != 0 {
		row := scoutRow + rowDirection
		for row != moveRow+rowDirection {
			if board.board[row][scoutCol] != nil {
				if board.board[row][scoutCol].Team == nil || *board.board[row][scoutCol].Team == scoutTeam {
					// cannot move over same team unit or water
					return false
				} else if *board.board[row][scoutCol].Team != scoutTeam && row != moveRow {
					// cannot move over enemy team unit unless last unit
					return false
				}
			}
			row = row + rowDirection
		}
	}
	colDirection := 0
	if moveCol-scoutCol > 0 {
		colDirection = 1
	} else if moveCol-scoutCol < 0 {
		colDirection = -1
	}
	if colDirection != 0 {
		col := scoutCol + colDirection
		for col != moveCol+colDirection {
			if board.board[scoutRow][col] != nil {
				if board.board[scoutRow][col].Team == nil || *board.board[scoutRow][col].Team == scoutTeam {
					// cannot move over same team unit or water
					return false
				} else if *board.board[scoutRow][col].Team != scoutTeam && col != moveCol {
					// cannot move over enemy team unit unless last unit
					return false
				}
			}
			col = col + colDirection
		}
	}
	return true
}

func (s *state) nextTurn() {
	if len(s.winners) > 0 {
		return
	}
	for idx, t := range s.teams {
		if t == s.turn {
			s.turn = s.teams[(idx+1)%len(s.teams)]
			return
		}
	}
}

func (s *state) playersReady() bool {
	r := true
	for _, b := range s.ready {
		r = r && b
	}
	return r
}

func (s *state) SetWinners(winners []string) error {
	for _, winner := range winners {
		if !contains(s.teams, winner) {
			return &bgerr.Error{
				Err:    fmt.Errorf("winner not in teams"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
	}
	s.winners = winners
	return nil
}

func (s *state) targets() []*bg.BoardGameAction {
	targets := make([]*bg.BoardGameAction, 0)
	for r, row := range s.board.board {
		for c, unit := range row {
			if unit != nil && unit.Team != nil &&
				*unit.Team == s.turn {
				for _, move := range s.board.possibleMoves(r, c) {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionMoveUnit,
						MoreDetails: MoveUnitActionDetails{
							UnitRow:    r,
							UnitColumn: c,
							MoveRow:    move[0],
							MoveColumn: move[1],
						},
					})
				}
			}
		}
	}
	return targets
}

func (s *state) message() string {
	var message string
	if s.started {
		message = fmt.Sprintf("%s must move a unit", s.turn)
	} else {
		message = "arrange your units"
	}
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s wins", s.winners[0])
	}
	return message
}
