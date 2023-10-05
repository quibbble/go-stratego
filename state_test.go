package go_stratego

import (
	"math/rand"
	"testing"
	"time"

	"github.com/quibbble/go-boardgame/pkg/bgn"
	"github.com/stretchr/testify/assert"
)

func Test_StateRandomness(t *testing.T) {
	seed := time.Now().Unix()
	teams := []string{"A", "B"}
	r1 := rand.New(rand.NewSource(seed))
	r2 := rand.New(rand.NewSource(seed))

	s1, _ := newState(teams, VariantQuickBattle, r1)
	s2, _ := newState(teams, VariantQuickBattle, r2)

	for i := 0; i < len(s1.board.board); i++ {
		for j := 0; j < len(s1.board.board); j++ {
			u1 := s1.board.board[i][j]
			u2 := s2.board.board[i][j]
			if u1 != nil {
				assert.Equal(t, u1.Team, u2.Team)
				assert.Equal(t, u1.Type, u2.Type)
			}
		}
	}
}

func Test_BGN(t *testing.T) {
	raw := `
		[Variant "QuickBattle"]
		[Game "Stratego"]
		[Teams "red, blue"]
		[Seed "1696470849747"]

		1s&6.5.5.4 1s&6.5.5.4
		0s&1.1.0.4 0s&0.4.0.5 0s&0.5.1.6 0s&1.6.1.4 0s&1.4.1.3 0s&1.3.0.3 0s&0.3.0.2
		0s&0.2.0.1 0s&0.1.0.0 0s&0.1.1.2 0s&1.2.1.4 0s&1.4.0.7 0s&1.6.0.6 0s&2.5.0.5
		0s&2.4.0.4 0s&1.3.0.1 1r 0r 0m&2.6.5.6`

	game, err := bgn.Parse(raw)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	builder := Builder{}
	if _, err := builder.Load(game); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
