package go_stratego

import (
	"fmt"
	"strconv"

	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{
		ActionSwitchUnits: "s",
		ActionToggleReady: "r",
		ActionMoveUnit:    "m",
	}
	notationToAction = reverseMap(actionToNotation)
)

func (s *SwitchUnitsActionDetails) encodeBGN() []string {
	return []string{strconv.Itoa(s.UnitRow), strconv.Itoa(s.UnitColumn), strconv.Itoa(s.SwitchUnitRow), strconv.Itoa(s.SwitchUnitColumn)}
}

func decodeSwitchUnitsActionDetailsBGN(notation []string) (*SwitchUnitsActionDetails, error) {
	if len(notation) != 4 {
		return nil, loadFailure(fmt.Errorf("invalid switch units notation"))
	}
	unitRow, err := strconv.Atoi(notation[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	unitColumn, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	switchUnitRow, err := strconv.Atoi(notation[2])
	if err != nil {
		return nil, loadFailure(err)
	}
	switchUnitColumn, err := strconv.Atoi(notation[3])
	if err != nil {
		return nil, loadFailure(err)
	}
	return &SwitchUnitsActionDetails{
		UnitRow:          unitRow,
		UnitColumn:       unitColumn,
		SwitchUnitRow:    switchUnitRow,
		SwitchUnitColumn: switchUnitColumn,
	}, nil
}

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

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
