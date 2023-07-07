package go_stratego

import (
	"fmt"
	"strconv"

	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{ActionMoveUnit: "m", ActionBattle: "b"}
	notationToAction = reverseMap(actionToNotation)
)

func (m *MoveUnitActionDetails) encodeBGN() []string {
	return []string{strconv.Itoa(m.UnitRow), strconv.Itoa(m.UnitColumn), strconv.Itoa(m.MoveRow), strconv.Itoa(m.MoveColumn)}
}

func decodeMoveUnitActionDetailsBGN(notation []string) (*MoveUnitActionDetails, error) {
	if len(notation) != 4 {
		return nil, loadFailure(fmt.Errorf("invalid move unit notation"))
	}
	unitRow, err := strconv.Atoi(notation[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	unitColumn, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	moveRow, err := strconv.Atoi(notation[2])
	if err != nil {
		return nil, loadFailure(err)
	}
	moveColumn, err := strconv.Atoi(notation[3])
	if err != nil {
		return nil, loadFailure(err)
	}
	return &MoveUnitActionDetails{
		UnitRow:    unitRow,
		UnitColumn: unitColumn,
		MoveRow:    moveRow,
		MoveColumn: moveColumn,
	}, nil
}

func (b *BattleActionDetails) encodeBGN() []string {
	return []string{b.AttackingUnit, b.AttackedUnit, b.Winner}
}

func decodeBattleActionDetailsBGN(notation []string) (*BattleActionDetails, error) {
	if len(notation) != 3 {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	return &BattleActionDetails{
		AttackingUnit: notation[0],
		AttackedUnit:  notation[1],
		Winner:        notation[2],
	}, nil
}

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
