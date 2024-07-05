using System;
using System.Collections.Generic;
using Xunit;

namespace BattleshipTest
{
    public class Coordinate
    {
        public int Num { get; set; }
        public char Letter { get; set; }
    }

    public class Position
    {
        public Coordinate Start { get; set; }

        public Coordinate End { get; set; }
    }

    public class IncorrectLetterException : Exception
    {
        public IncorrectLetterException()
        {
        }

        public IncorrectLetterException(string message)
            : base(message)
        {
        }

        public IncorrectLetterException(string message, Exception inner)
            : base(message, inner)
        {
        }
    }

    public class OutOfGridBoundariesException : Exception
    {
        public OutOfGridBoundariesException()
        {
        }

        public OutOfGridBoundariesException(string message)
            : base(message)
        {
        }

        public OutOfGridBoundariesException(string message, Exception inner)
            : base(message, inner)
        {
        }
    }

    public class Grid
    {
        public class Ship
        {
            public Coordinate Start { get; set; }

            public Coordinate End { get; set; }
        }

        public Position GridSize => new Position
        {
            Start = new Coordinate { Num = 1, Letter = 'A' },
            End = new Coordinate { Num = 10, Letter = 'J' },
        };

        private Ship[] _ships { get; set; }

        public int Shots { get; private set; }

        public Grid(Position[] ships)
        {
            _ships = new Ship[ships.Length];

            for (var i = 0; i < ships.Length; i++)
            {
                var ship = new Ship
                {
                    Start = ships[i].Start,
                    End = ships[i].End,
                };
                _ships[i] = ship;
            }
        }

        public ShootResult Shoot(int shotNum, string shotLetter)
        {
            //TODO: implement here
            return new ShootResult();
        }

        public void ResetGrid()
        {
            Shots = 0;
        }
    }

    public class ShootResult
    {
        public bool Hit { get; set; }
        public bool Sunk { get; set; }
    }
	
    public class GridTests
    {
        public class TShot
        {
            public int Num { get; set; }
            public string Letter { get; set; }
            public bool ExpectedHit { get; set; }
        }

        [Fact]
        public void TestMissHitMissNotSunk()
        {
            var shots = new TShot[]
            {
                new TShot { Num = 1, Letter = "G", ExpectedHit = false },
                new TShot { Num = 1, Letter = "H", ExpectedHit = true },
                new TShot { Num = 1, Letter = "I", ExpectedHit = false },
            };

            TestShoots(shots, false);
        }

        [Fact]
        public void TestSinkBattleship()
        {
            var shots = new TShot[]
            {
                new TShot { Num = 1, Letter = "H", ExpectedHit = true },
                new TShot { Num = 2, Letter = "H", ExpectedHit = true },
                new TShot { Num = 3, Letter = "H", ExpectedHit = true },
                new TShot { Num = 4, Letter = "H", ExpectedHit = true }
            };

            TestShoots(shots, true);
        }

        [Fact]
        public void TestMissThenSink()
        {
            var shots = new TShot[]
            {
                new TShot { Num = 1, Letter = "D", ExpectedHit = false },
                new TShot { Num = 7, Letter = "F", ExpectedHit = true },
                new TShot { Num = 8, Letter = "F", ExpectedHit = true }
            };

            TestShoots(shots, true);
        }

        [Fact]
        public void TestMissFewTimesAndSink()
        {
            var shots = new TShot[]
            {
                new TShot { Num = 10, Letter = "D", ExpectedHit = true },
                new TShot { Num = 9, Letter = "D", ExpectedHit = false },
                new TShot { Num = 10, Letter = "C", ExpectedHit = false },
                new TShot { Num = 10, Letter = "E", ExpectedHit = true },
                new TShot { Num = 10, Letter = "F", ExpectedHit = true },
                new TShot { Num = 10, Letter = "G", ExpectedHit = true },
                new TShot { Num = 10, Letter = "H", ExpectedHit = true }
            };

            TestShoots(shots, true);
        }


        [Fact]
        public void TestDiagonalShots()
        {
            var shots = new TShot[] 
            { 
                new TShot { Num = 7, Letter = "J", ExpectedHit = true },
                new TShot { Num = 8, Letter = "I", ExpectedHit = false },
                new TShot { Num = 9, Letter = "H", ExpectedHit = false },
                new TShot { Num = 10, Letter = "G", ExpectedHit = true } 
            };

            TestShoots(shots, true);
        }

        private void TestShoots(TShot[] shots, bool expectedSunk)
        {
            var ships = GetShips();

            var grid = new Grid(ships);

            var latestShootResult = new ShootResult();

            foreach (var shot in shots)
            {

                latestShootResult = grid.Shoot(shot.Num, shot.Letter);

                Assert.True(shot.ExpectedHit == latestShootResult.Hit, "hit result is not as expected");
            }

            Assert.True(expectedSunk == latestShootResult.Sunk, "sunk result is not as expected");
        }

          [Fact]
        public void TestShootThrowsIncorrectLetterException()
        {
            var ships = GetShips();

            var grid = new Grid(ships);
            Assert.Throws<IncorrectLetterException>(() => grid.Shoot(10, "CC"));

        }

        [Fact]
        public void TestShootThrowsOutOfGridBoundariesException_Num()
        {
            var ships = GetShips();

            var grid = new Grid(ships);
            Assert.Throws<OutOfGridBoundariesException>(() => grid.Shoot(11, "G"));
        }

        [Fact]
        public void TestShootThrowsOutOfGridBoundariesException_Letter()
        {
            var ships = GetShips();

            var grid = new Grid(ships);
            Assert.Throws<OutOfGridBoundariesException>(() => grid.Shoot(10, "P"));
        }

        [Fact]
        public void TestShootThrowsOutOfGridBoundariesException_Both<T>()
        {
            var ships = GetShips();

            var grid = new Grid(ships);
            Assert.Throws<OutOfGridBoundariesException>(() => grid.Shoot(12, "K"));
        }

        private Position[] GetShips()
        {
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
            //		7           @     @
            //		8           @     @
            //		9                 @
            //	   10       @ @ @ @ @
            //
            var ships = new List<Position>();
            ships.Add(new Position
            {
                Start = new Coordinate { Num = 2, Letter = 'A' },
                End = new Coordinate { Num = 2, Letter = 'A' },
            });

            ships.Add(new Position
            {
                Start = new Coordinate { Num = 3, Letter = 'E' },
                End = new Coordinate { Num = 3, Letter = 'E' },
            });
            ships.Add(new Position
            {
                Start = new Coordinate { Num = 1, Letter = 'H' },
                End = new Coordinate { Num = 4, Letter = 'H' },
            });
            ships.Add(new Position
            {
                Start = new Coordinate { Num = 5, Letter = 'B' },
                End = new Coordinate { Num = 5, Letter = 'C' },
            });

            ships.Add(new Position
            {
                Start = new Coordinate { Num = 7, Letter = 'F' },
                End = new Coordinate { Num = 8, Letter = 'F' },
            });

            ships.Add(new Position
            {
                Start = new Coordinate { Num = 7, Letter = 'I' },
                End = new Coordinate { Num = 9, Letter = 'I' },
            });

            ships.Add(new Position
            {
                Start = new Coordinate { Num = 10, Letter = 'D' },
                End = new Coordinate { Num = 10, Letter = 'H' },
            });

            return ships.ToArray();
        }
    }
}
