package tests_test

import (
	"fmt"
	"testing"

	tests "github.com/SnowSoftwareGlobal/saas-codetests"
	"github.com/stretchr/testify/assert"
)

type shot struct {
	num         int
	letter      string
	expectedHit bool
}

type testCase struct {
	shots        []shot
	expectedSunk bool
	err          error
}

func TestShoot(t *testing.T) {
	ships := getShips()

	testCases := []testCase{
		{shots: []shot{{1, "G", false}, {1, "H", true}, {1, "I", false}}, expectedSunk: false},
		{shots: []shot{{1, "H", true}, {2, "H", true}, {3, "H", true}, {4, "H", true}}, expectedSunk: true},
		{shots: []shot{{1, "D", false}, {7, "F", true}, {8, "F", true}}, expectedSunk: true},
		{shots: []shot{{10, "D", true}, {9, "D", false}, {10, "C", false},
			{10, "E", true}, {10, "F", true}, {10, "G", true}, {10, "H", true}},
			expectedSunk: true},
		{shots: []shot{{7, "J", true}, {8, "I", false}, {9, "H", false}, {10, "G", true}}, expectedSunk: false},
		{shots: []shot{{1, "H", true}, {8, "I", false}, {10, "CC", false}, {10, "G", true}}, expectedSunk: false, err: tests.ErrIncorrectLetter},
		{shots: []shot{{1, "G", false}, {1, "H", true}, {11, "G", false}}, expectedSunk: false, err: tests.ErrOutOfGridBoundaries},
		{shots: []shot{{1, "G", false}, {1, "H", true}, {10, "P", false}}, expectedSunk: false, err: tests.ErrOutOfGridBoundaries},
		{shots: []shot{{1, "G", false}, {12, "K", false}, {11, "G", false}}, expectedSunk: false, err: tests.ErrOutOfGridBoundaries},
	}
	for ind, test := range testCases {
		t.Run(fmt.Sprint(ind), func(t *testing.T) {
			grid, err := tests.NewGrid(ships)
			assert.NoErrorf(t, err, "failed to create grid")
			var latestShootResult tests.ShootResult
			for _, shot := range test.shots {
				var err error
				latestShootResult, err = grid.Shoot(shot.num, shot.letter)
				assert.Equal(t, shot.expectedHit, latestShootResult.Hit)
				if err != nil && test.err != nil {
					assert.ErrorIsf(t, test.err, err, "err is incorrect")
				}
			}

			assert.Equal(t, test.expectedSunk, latestShootResult.Sunk)

			// you do this when you want to play again the same game
			grid.ResetShips()
			// you do this when you don't want to play the same game again
			// this will help out GC to clean up faster
			grid.Destroy()
		})
	}
}

func getShips() []tests.Position {
	// count  name              size
	//   1    Aircraft Carrier   5
	//   1    Battleship         4
	//   1    Cruiser            3
	//   2    Destroyer          2
	//   2    Submarine          1
	//
	// 		  A B C D E F G H I J
	//		1               @
	//		2 @             @
	//		3         @     @
	//		4               @
	//		5   @ @
	//		6
	//		7           @       @
	//		8           @       @
	//		9                   @
	//	   10       @ @ @ @ @
	//
	var ships []tests.Position
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{2, 'A'},
		End:   tests.Coordinate{2, 'A'},
	})
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{3, 'E'},
		End:   tests.Coordinate{3, 'E'},
	})
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{1, 'H'},
		End:   tests.Coordinate{4, 'H'},
	})
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{5, 'B'},
		End:   tests.Coordinate{5, 'C'},
	})
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{7, 'F'},
		End:   tests.Coordinate{8, 'F'},
	})
	// there was an error here letter was "I" instead of "J"
	// (which is shown in the picture)
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{7, 'J'},
		End:   tests.Coordinate{9, 'J'},
	})
	ships = append(ships, tests.Position{
		Start: tests.Coordinate{10, 'D'},
		End:   tests.Coordinate{10, 'H'},
	})

	return ships
}
