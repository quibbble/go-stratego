package go_stratego

import "fmt"

// UnitTyes ordered in ascending order by battle winner
var UnitTyes = []string{
	"flag", "bomb", "spy", "scout", "miner", "sergeant", "lieutenant",
	"captain", "major", "colonel", "general", "marshal",
}

type Unit struct {
	typ  string
	team *string
}

func NewUnit(typ, team string) *Unit {
	return &Unit{
		typ:  typ,
		team: &team,
	}
}

func Water() *Unit {
	return &Unit{
		typ:  "water",
		team: nil,
	}
}

func (u *Unit) Attack(unit *Unit) (winner *Unit, err error) {
	if *u.team == *unit.team {
		return nil, fmt.Errorf("cannot attack unit on same team")
	}
	if u.typ == "flag" || u.typ == "bomb" {
		return nil, fmt.Errorf("%s cannot attack", u.typ)
	}
	// spy -> marshal case
	if u.typ == "spy" && unit.typ == "marshal" {
		return u, nil
	}
	// miner -> bomb case
	if u.typ == "miner" && unit.typ == "bomb" {
		return u, nil
	}
	// any -> bomb case
	if unit.typ == "bomb" {
		return nil, nil
	}
	// same type case
	if u.typ == unit.typ {
		return nil, nil
	}
	// default case
	if indexOf(UnitTyes, u.typ) > indexOf(UnitTyes, unit.typ) {
		return u, nil
	} else {
		return unit, nil
	}
}
