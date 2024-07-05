package tests

import (
	"errors"
)

type Coordinate struct {
	Num    int
	Letter rune
}

// Used for placing ships on the grid
type Position struct {
	Start Coordinate
	End   Coordinate
}

var (
	// added some errors to cover additional cases
	ErrIncorrectLetter         = errors.New("incorrect input in string")
	ErrOutOfGridBoundaries     = errors.New("shot out of grid boundaries")
	ErrShipOutOfGridBoundaries = errors.New("ship out of grid boundaries")
	ErrIncorrectPlacement      = errors.New("ship placed incorrectly")
	ErrShipsOverlap            = errors.New("two ships occupy same coordinate")

	TotalMiss = ShootResult{Hit: false, Sunk: false}
)

// removed GridSize as it is not needed any more

// I was considering just extending Cooridinate struct but decided against it.
//
// This seems more natural. I can avoid circular pointing.
type gridPoint struct {
	hasBeenHit bool
	shipPntr   *ship
}

func (gp *gridPoint) takeShot() ShootResult {
	// check if grip point was hit before
	if !gp.hasBeenHit {
		gp.hasBeenHit = true
		gp.shipPntr.hits++
	}
	return ShootResult{
		Hit:  true,
		Sunk: gp.shipPntr.hits == gp.shipPntr.gridLength, // determine if the ship is sunk
	}
}

func (gp *gridPoint) reset() {
	gp.hasBeenHit = false
	gp.shipPntr.hits = 0
}

func (gp *gridPoint) destroy() {
	gp.shipPntr = nil
}

// Map of grid points. Key is the letter partof the shot coordinate.
// The value is the grid point being hit.
type lettersMap map[rune]gridPoint

func (lm lettersMap) addShip(letter rune, ship *ship) error {
	_, ok := lm[letter]
	if ok {
		return ErrShipsOverlap
	}
	lm[letter] = gridPoint{
		hasBeenHit: false,
		shipPntr:   ship,
	}
	return nil
}

func (lm lettersMap) shoot(letter rune) ShootResult {
	gp, ok := lm[letter]
	if !ok {
		return TotalMiss
	}
	return gp.takeShot()
}

func (lm lettersMap) reset() {
	for _, gp := range lm {
		gp.reset()
	}
}

func (lm lettersMap) destroy() {
	for _, gp := range lm {
		gp.destroy()
	}
}

// Key in this map is the number part of the shot coordinate.
//
// I realize I could've organized these maps in some generic
// struct but I think it's overkill for something this simple
type numbersMap map[int]lettersMap

func (nm numbersMap) addShip(number int, letter rune, ship *ship) error {
	lm, ok := nm[number]
	if !ok {
		lm = make(lettersMap)
		nm[number] = lm
	}
	return lm.addShip(letter, ship)
}

func (nm numbersMap) shoot(number int, letter rune) ShootResult {
	lm, ok := nm[number]
	if ok {
		return lm.shoot(letter)
	}
	return TotalMiss
}

func (nm numbersMap) reset() {
	for _, lm := range nm {
		lm.reset()
	}
}

func (nm numbersMap) destroy() {
	for key, lm := range nm {
		lm.destroy()
		nm[key] = nil
	}
}

type Grid struct {
	// I'm leaving this in order to be able to iterate over ships
	// but right now there is no need for it
	ships []*ship
	// Added this structure to speed up hit/miss resolution
	shipMap numbersMap
}

type ship struct {
	// how many times a ship has been uniquely hit
	hits int
	// how many unique hits a ship can take before being sunk
	gridLength int
	start      Coordinate
	end        Coordinate
}

type ShootResult struct {
	Hit  bool
	Sunk bool
}

// Make sure coordinates are within bounds of the grid
func coordinatesAreValid(number int, letter rune) bool {
	return number >= 0 || number <= 10 || letter >= 'A' || letter <= 'J'
}

// Adds all coordinates of the ship to the ships map.
// This will speed up shot hit/miss resolution later on.
func addShipToGridAndSetShipLength(ship *ship, ships numbersMap) error {
	var (
		err error
		// counter ranges for adding ships
		iStart int
		iEnd   int
		jStart rune
		jEnd   rune
	)
	if ship.start.Num < ship.end.Num {
		iStart = ship.start.Num
		iEnd = ship.end.Num
	} else {
		iStart = ship.end.Num
		iEnd = ship.start.Num
	}
	if ship.start.Letter < ship.end.Letter {
		jStart = ship.start.Letter
		jEnd = ship.end.Letter
	} else {
		jStart = ship.end.Letter
		jEnd = ship.start.Letter
	}
	for i := iStart; i <= iEnd; i++ {
		for j := jStart; j <= jEnd; j++ {
			err = ships.addShip(i, j, ship)
			if err != nil {
				return err
			}
		}
	}
	// set ship length
	switch {
	case iEnd > iStart:
		ship.gridLength = iEnd - iStart + 1
	case jEnd > jStart:
		ship.gridLength = int(jEnd-jStart) + 1
	default:
		ship.gridLength = 1
	}
	return nil
}

func NewGrid(ships []Position) (*Grid, error) {
	var (
		realShips []*ship
		err       error
	)
	shipMap := make(numbersMap)
	for _, p := range ships {
		ship := ship{
			start: p.Start,
			end:   p.End,
		}
		// validate coordinates
		if !coordinatesAreValid(ship.start.Num, ship.start.Letter) || !coordinatesAreValid(ship.end.Num, ship.end.Letter) {
			return nil, ErrShipOutOfGridBoundaries
		}
		// make sure ship is place horizontally or vertically and not diagonally
		if ship.start.Num != ship.end.Num && ship.start.Letter != ship.end.Letter {
			return nil, ErrIncorrectPlacement
		}
		// add ship to grid
		if err = addShipToGridAndSetShipLength(&ship, shipMap); err != nil {
			return nil, err
		}
		// add ship to the list
		realShips = append(realShips, &ship)
	}

	return &Grid{
		ships:   realShips,
		shipMap: shipMap,
	}, nil
}

func (grid *Grid) Shoot(shotNum int, shotLetter string) (ShootResult, error) {
	if len(shotLetter) != 1 {
		return TotalMiss, ErrIncorrectLetter
	}
	letter := []rune(shotLetter)[0]
	if !coordinatesAreValid(shotNum, letter) {
		return TotalMiss, ErrOutOfGridBoundaries
	}
	return grid.shipMap.shoot(shotNum, letter), nil
}

// Reset all hit data so we can play the same again
func (grid *Grid) ResetShips() {
	if grid.shipMap != nil {
		grid.shipMap.reset()
	}
}

// Tear down entire grid structure to make things easier on the GC.
//
// This is in no way mandatory but I do it if the structure in question
// is compicated enough.
func (grid *Grid) Destroy() {
	if grid.shipMap != nil {
		grid.shipMap.destroy()
	}
}
