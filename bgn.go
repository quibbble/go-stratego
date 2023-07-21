package go_stratego

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{
		ActionSwitchUnits: "s",
		ActionMoveUnit:    "m",
		ActionBattle:      "b",
	}
	notationToAction   = reverseMap(actionToNotation)
	unitTypeToNotation = map[string]string{
		flag:       "f",
		bomb:       "b",
		spy:        "s",
		scout:      "1",
		miner:      "2",
		sergeant:   "3",
		lieutenant: "4",
		captain:    "5",
		major:      "6",
		colonel:    "7",
		general:    "8",
		marshal:    "9",
	}
	notationToUnitType = reverseMap(unitTypeToNotation)
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

func (b *BattleActionDetails) encodeBGN(teams []string) []string {
	attacking := fmt.Sprintf("%d:%s", indexOf(teams, *b.AttackingUnit.Team), unitTypeToNotation[b.AttackingUnit.Type])
	attacked := fmt.Sprintf("%d:%s", indexOf(teams, *b.AttackedUnit.Team), unitTypeToNotation[b.AttackedUnit.Type])
	return []string{attacking, attacked, strconv.Itoa(indexOf(teams, b.WinningTeam))}
}

func decodeBattleActionDetailsBGN(teams []string, notation []string) (*BattleActionDetails, error) {
	if len(notation) != 3 {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	attacking := strings.Split(notation[0], ":")
	attacked := strings.Split(notation[1], ":")
	if len(attacking) != 2 || len(attacked) != 2 {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	attackingIndex, err := strconv.Atoi(attacking[0])
	if err != nil || attackingIndex < 0 || attackingIndex >= len(teams) {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	attackedIndex, err := strconv.Atoi(attacked[0])
	if err != nil || attackedIndex < 0 || attackedIndex >= len(teams) {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	attackingType := notationToUnitType[attacking[1]]
	attackedType := notationToUnitType[attacked[1]]
	if attackingType == "" || attackedType == "" {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	winningIndex, err := strconv.Atoi(notation[2])
	if err != nil || winningIndex < -1 || winningIndex >= len(teams) {
		return nil, loadFailure(fmt.Errorf("invalid battle notation"))
	}
	winningTeam := ""
	if winningIndex >= 0 {
		winningTeam = teams[winningIndex]
	}
	return &BattleActionDetails{
		AttackingUnit: *NewUnit(teams[attackingIndex], attackingType),
		AttackedUnit:  *NewUnit(teams[attackedIndex], attackedType),
		WinningTeam:   winningTeam,
	}, nil
}

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
