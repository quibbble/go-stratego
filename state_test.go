package go_stratego

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_StateRandomness(t *testing.T) {
	seed := time.Now().Unix()
	teams := []string{"A", "B"}
	r1 := rand.New(rand.NewSource(seed))
	r2 := rand.New(rand.NewSource(seed))

	s1, _ := newState(teams, r1)
	s2, _ := newState(teams, r2)

	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			u1 := s1.board.board[i][j]
			u2 := s2.board.board[i][j]
			if u1 != nil {
				assert.Equal(t, u1.Team, u2.Team)
				assert.Equal(t, u1.Type, u2.Type)
			}

		}
	}
}
