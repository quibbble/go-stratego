package go_stratego

const (
	BoardSize            = 10
	QuickBattleBoardSize = 8
)

// Action types
const (
	ActionSwitchUnits = "SwitchUnits"
	ActionMoveUnit    = "MoveUnit"
	ActionToggleReady = "ToggleReady"
)

// Stratego Variants
const (
	VariantClassic     = "Classic"     // normal Stratego
	VariantQuickBattle = "QuickBattle" // 8x8 quick play Stratego
)

var Variants = []string{VariantClassic, VariantQuickBattle}

// StrategoMoreOptions are the additional options for creating a game of Stratego
type StrategoMoreOptions struct {
	Seed    int64
	Variant string
}

// StategoSnapshotData is the game data unique to Statego
type StategoSnapshotData struct {
	Board       [][]Unit
	Battle      *Battle
	JustBattled bool
	Started     bool
	Ready       map[string]bool
	Variant     string
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
