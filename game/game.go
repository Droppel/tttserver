package game

import (
	"fmt"
)

type Point struct {
	x, y int
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%d|%d", p.x, p.y)
}

func AddPoints(p1, p2 Point) Point {
	p := InitPoint(p1.x+p2.x, p1.y+p2.y)
	return p
}

func SubPoints(p1, p2 Point) Point {
	p := InitPoint(p1.x-p2.x, p1.y-p2.y)
	return p
}

type Game struct {
	size    int
	winSize int
	turn    int
	state   map[Point]int
}

func InitPoint(x, y int) Point {
	return Point{x, y}
}

func InitGame(size, winSize int) *Game {
	return &Game{
		size:    size,
		winSize: winSize,
		state:   make(map[Point]int),
		turn:    0,
	}
}

func (g *Game) MakeMove(move Point) bool {
	_, exists := g.state[move]
	if exists {
		return false
	}

	if !g.CheckInbound(move) {
		return false
	}

	g.state[move] = g.GetActivePlayerNum()
	g.turn++
	return true
}

func (g *Game) GetActivePlayerNum() int {
	return (g.turn)%2 + 1
}

func (g *Game) StateToString() string {
	out := fmt.Sprintf("STATE ")
	for k, v := range g.state {
		out += fmt.Sprintf("%s|%d,", k.ToString(), v)
	}
	return out
}

var (
	HORIZONTAL    = InitPoint(1, 0)
	VERTICAL      = InitPoint(0, 1)
	DIAGONAL_UP   = InitPoint(1, 1)
	DIAGONAL_DOWN = InitPoint(1, -1)
	DIRECTIONS    = []Point{HORIZONTAL, VERTICAL, DIAGONAL_UP, DIAGONAL_DOWN}
)

func (g *Game) CheckDone(move Point) (bool, int) {
	currentSymbol := g.state[move]

	for _, dir := range DIRECTIONS {
		totalCount := 1
		currentPos := InitPoint(move.x, move.y)
		for {
			currentPos = AddPoints(currentPos, dir)
			if !g.CheckInbound(currentPos) || g.state[currentPos] != currentSymbol {
				break
			}
			totalCount++
		}
		currentPos = InitPoint(move.x, move.y)
		for {
			currentPos = SubPoints(currentPos, dir)
			if !g.CheckInbound(currentPos) || g.state[currentPos] != currentSymbol {
				break
			}
			totalCount++
		}
		if totalCount >= g.winSize {
			return true, currentSymbol
		}
	}

	//Check draw
	if len(g.state) >= g.size*g.size {
		return true, 0
	}
	return false, 0
}

func (g *Game) CheckInbound(p Point) bool {
	maxVal := (g.size - 1) / 2
	minVal := -1 * maxVal
	return !(p.x > maxVal || p.x < minVal || p.y > maxVal || p.y < minVal)
}
