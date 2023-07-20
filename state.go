package go_stratego

import (
	"fmt"
	"math"
	"math/rand"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

type state struct {
	turn    string
	teams   []string
	winners []string
	board   *Board
	started bool
}

func newState(teams []string, random *rand.Rand) (*state, error) {
	if random == nil {
		return nil, fmt.Errorf("random seed is null")
	}
	board, err := NewRandomBoard(teams, random)
	if err != nil {
		return nil, err
	}
	return &state{
		board:   board,
		teams:   teams,
		turn:    teams[0],
		winners: make([]string, 0),
	}, nil
}

func (s *state) SwitchUnits(team string, unitRow, unitCol, switchUnitRow, switchUnitCol int) error {
	if s.started {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units when game has already started"),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if unitRow >= BoardSize || unitRow < 0 || unitCol >= BoardSize || unitCol < 0 ||
		switchUnitRow >= BoardSize || switchUnitRow < 0 || switchUnitCol >= BoardSize || switchUnitCol < 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("index out of bounds"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	unit := s.board.board[unitRow][unitCol]
	switchUnit := s.board.board[switchUnitRow][switchUnitCol]
	if unit.Team == nil || switchUnit.Team == nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot switch units that have no team"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if *unit.Team != *switchUnit.Team {
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
	s.board.board[unitRow][unitCol] = switchUnit
	s.board.board[switchUnitRow][switchUnitCol] = unit
	return nil
}

func (s *state) MoveUnit(team string, unitRow, unitCol, moveRow, moveCol int) (*BattleActionDetails, error) {
	s.started = true
	if team != s.turn {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("%s cannot play on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if unitRow >= BoardSize || unitRow < 0 || unitCol >= BoardSize || unitCol < 0 ||
		moveRow >= BoardSize || moveRow < 0 || moveCol >= BoardSize || moveCol < 0 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("index out of bounds"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	unit := s.board.board[unitRow][unitCol]
	if unit == nil {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("unit does not exist at %d, %d", unitRow, unitCol),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if *unit.Team != team {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("cannot move a unit not part of your team"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if unit.Type == "bomb" || unit.Type == "flag" {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("cannot move bombs or flags"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if math.Abs(float64(moveRow)-float64(unitRow))+math.Abs(float64(moveCol)-float64(unitCol)) > 1.0 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("unit cannot move diagonally"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if (math.Abs(float64(moveRow)-float64(unitRow)) > 1.0 && unit.Type != "scout") ||
		(math.Abs(float64(moveCol)-float64(unitCol)) > 1.0 && unit.Type != "scout") {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("unit cannot move more than one space unless they are a scout"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if unit.Type == "scout" && !scoutCanMove(s.board, unitRow, unitCol, moveRow, moveCol, *unit.Team) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("scout cannot move through water or other units"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	attackedUnit := s.board.board[moveRow][moveCol]
	if attackedUnit != nil {
		if attackedUnit.Type == "water" {
			return nil, &bgerr.Error{
				Err:    fmt.Errorf("cannot move onto water"),
				Status: bgerr.StatusInvalidAction,
			}
		}
		winningUnit, err := unit.Attack(attackedUnit)
		if err != nil {
			return nil, err
		}
		s.board.board[unitRow][unitCol] = nil
		s.board.board[moveRow][moveCol] = winningUnit
		winner := ""
		if winningUnit != nil {
			winner = *winningUnit.Team
		}
		// check for game over
		if attackedUnit.Type == "flag" {
			s.winners = []string{team} // attacked flag so game is over
		} else if s.board.numActive(*attackedUnit.Team) == 0 && s.board.numActive(team) == 0 {
			s.winners = []string{""} // both teams ran out of movable units
		} else if s.board.numActive(*attackedUnit.Team) == 0 {
			s.winners = []string{team} // one team ran out of movable units
		} else if s.board.numActive(team) == 0 {
			s.winners = []string{*attackedUnit.Team} // the other team ran out of movable units
		}
		s.nextTurn()
		return &BattleActionDetails{
			AttackingUnit: *unit,
			AttackedUnit:  *attackedUnit,
			WinningTeam:   winner,
		}, nil
	} else {
		s.board.board[unitRow][unitCol] = nil
		s.board.board[moveRow][moveCol] = unit
		s.nextTurn()
		return nil, nil
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
	message := fmt.Sprintf("%s must move a unit", s.turn)
	return message
}
