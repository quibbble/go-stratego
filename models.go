package go_stratego

const BoardSize = 10

// Action types
const (
	ActionMoveUnit = "MoveUnit"
	ActionBattle   = "Battle"
)

type MoveUnitActionDetails struct {
	UnitRow, UnitColumn int
	MoveRow, MoveColumn int
}

type BattleActionDetails struct {
	AttackingUnit, AttackedUnit, Winner string
}
