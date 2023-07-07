package go_stratego

const BoardSize = 10

// Action types
const (
	ActionMoveUnit = "MoveUnit"
	ActionBattle   = "Battle"
)

// StrategoMoreOptions are the additional options for creating a game of Stratego
type StrategoMoreOptions struct {
	Seed int64
}

type MoveUnitActionDetails struct {
	UnitRow, UnitColumn int
	MoveRow, MoveColumn int
}

type BattleActionDetails struct {
	AttackingUnit, AttackedUnit, Winner string
}
