package go_stratego

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
)

const (
	minTeams = 2
	maxTeams = 2
)

type Stratego struct {
	state   *state
	actions []*bg.BoardGameAction
	options *StrategoMoreOptions
}

func NewStratego(options *bg.BoardGameOptions) (*Stratego, error) {
	if len(options.Teams) < minTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key),
			Status: bgerr.StatusTooFewTeams,
		}
	} else if len(options.Teams) > maxTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key),
			Status: bgerr.StatusTooManyTeams,
		}
	} else if duplicates(options.Teams) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("duplicate teams found"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	var details StrategoMoreOptions
	if err := mapstructure.Decode(options.MoreOptions, &details); err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	state, err := newState(options.Teams, rand.New(rand.NewSource(details.Seed)))
	if err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	return &Stratego{
		state:   state,
		actions: make([]*bg.BoardGameAction, 0),
	}, nil
}

func (s *Stratego) Do(action *bg.BoardGameAction) error {
	if len(s.state.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("game already over"),
			Status: bgerr.StatusGameOver,
		}
	}
	switch action.ActionType {
	case ActionMoveUnit:
		var details MoveUnitActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		battle, err := s.state.MoveUnit(action.Team, details.UnitRow, details.UnitColumn, details.MoveRow, details.MoveColumn)
		if err != nil {
			return err
		}
		s.actions = append(s.actions, action, &bg.BoardGameAction{
			Team:        action.Team,
			ActionType:  ActionBattle,
			MoreDetails: battle,
		})
	case ActionBattle:
		// do nothing - this is here so BGN processing does not fail
	case bg.ActionSetWinners:
		var details bg.SetWinnersActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := s.state.SetWinners(details.Winners); err != nil {
			return err
		}
		s.actions = append(s.actions, action)
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (s *Stratego) GetSnapshot(team ...string) (*bg.BoardGameSnapshot, error) {
	if len(team) > 1 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("get snapshot requires zero or one team"),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	var targets = []*bg.BoardGameAction{}
	if len(s.state.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == s.state.turn)) {
		targets = s.state.targets()
	}

	board := [BoardSize][BoardSize]*Unit{}
	for r, row := range s.state.board.board {
		for c, unit := range row {
			if unit != nil {
				if unit.team != nil {
					if len(team) == 1 {
						if *unit.team == team[0] {
							board[r][c] = NewUnit(unit.typ, *unit.team)
						} else {
							board[r][c] = NewUnit("", *unit.team)
						}
					} else {
						board[r][c] = NewUnit("", *unit.team)
					}
				} else {
					board[r][c] = Water()
				}
			}
		}
	}

	return &bg.BoardGameSnapshot{
		Turn:    s.state.turn,
		Teams:   s.state.teams,
		Winners: s.state.winners,
		MoreData: StategoSnapshotData{
			Board: board,
		},
		Targets: targets,
		Actions: s.actions,
		Message: s.state.message(),
	}, nil
}

func (s *Stratego) GetBGN() *bgn.Game {
	tags := map[string]string{
		"Game":  key,
		"Teams": strings.Join(s.state.teams, ", "),
		"Seed":  fmt.Sprintf("%d", s.options.Seed),
	}
	actions := make([]bgn.Action, 0)
	for _, action := range s.actions {
		bgnAction := bgn.Action{
			TeamIndex: indexOf(s.state.teams, action.Team),
			ActionKey: rune(actionToNotation[action.ActionType][0]),
		}
		switch action.ActionType {
		case ActionMoveUnit:
			var details MoveUnitActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case ActionBattle:
			var details BattleActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case bg.ActionSetWinners:
			var details bg.SetWinnersActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details, _ = details.EncodeBGN(s.state.teams)
		}
		actions = append(actions, bgnAction)
	}
	return &bgn.Game{
		Tags:    tags,
		Actions: actions,
	}
}
