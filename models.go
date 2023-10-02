package go_stratego

const BoardSize = 10

// Action types
const (
	ActionSwitchUnits = "SwitchUnits"
	ActionMoveUnit    = "MoveUnit"
)

// StrategoMoreOptions are the additional options for creating a game of Stratego
type StrategoMoreOptions struct {
	Seed int64
}

// StategoSnapshotData is the game data unique to Statego
type StategoSnapshotData struct {
	Board       [BoardSize][BoardSize]*Unit
	Battle      *Battle
	JustBattled bool
	Started     bool
}

type SwitchUnitsActionDetails struct {
	UnitRow, UnitColumn             int
	SwitchUnitRow, SwitchUnitColumn int
}

type MoveUnitActionDetails struct {
	UnitRow, UnitColumn int
	MoveRow, MoveColumn int
}

type Battle struct {
	MoveUnitActionDetails
	AttackingUnit, AttackedUnit Unit
	WinningTeam                 string
}
