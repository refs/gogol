package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	aliveCell = '\u2B1B'
	deadCell  = '\u2B1C'
)

type board struct {
	Width  int
	Height int
	Cells  map[int][]Cell // encode locations of Alive Cells
	Clock  time.Duration  // how often a cycle lasts
}

func (b *board) runCycle() {
	cells := map[int][]Cell{}
	for i := 0; i < b.Height; i++ {
		cells[i] = make([]Cell, b.Width)
		for j := 0; j < b.Width; j++ {
			cells[i][j] = Cell{
				Alive: b.Cells[i][j].Alive,
				C: coordinates{
					x: b.Cells[i][j].C.x,
					y: b.Cells[i][j].C.y,
				},
			}
		}
	}

	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			aliveNeighbors := b.Cells[i][j].aliveNeighbors(*b)
			if b.Cells[i][j].Alive {
				if aliveNeighbors < 2 {
					cells[i][j].die() // underpopulation
				} else if aliveNeighbors == 2 || aliveNeighbors == 3 {
					continue // remain alive
				} else if aliveNeighbors > 3 {
					cells[i][j].die() // overpopulation
				}
				continue
			}

			if aliveNeighbors == 3 {
				cells[i][j].born()
			}
		}
	}

	b.Cells = cells
}

func newBoard(width, height int) board {
	b := board{
		Width:  width,
		Height: height,
		Cells:  map[int][]Cell{},
	}

	for i := 0; i < height; i++ {
		b.Cells[i] = make([]Cell, width)
		for j := 0; j < width; j++ {
			b.Cells[i][j].Alive = false
			b.Cells[i][j].C.x = i
			b.Cells[i][j].C.y = j
		}
	}
	return b
}

type Cell struct {
	Alive bool
	C     coordinates
}

// neighbors returns all the neighbors of a Cell for on a board b
func (c *Cell) neighbors(b board) []coordinates {
	neighbors := make([]coordinates, 0)
	// top left
	if !(c.C.y == 0) && !(c.C.x == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x - 1,
			y: c.C.y - 1,
		})
	}

	// top neighbor
	if !(c.C.y == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x,
			y: c.C.y - 1,
		})
	}

	// top right
	if !(c.C.x == 0) && !(c.C.y == b.Width-1) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x - 1,
			y: c.C.y + 1,
		})
	}

	// right
	if !(c.C.y == b.Width-1) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x,
			y: c.C.y + 1,
		})
	}

	// bottom right
	if !(c.C.y == b.Width-1) && !(c.C.x == b.Height-1) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x + 1,
			y: c.C.y + 1,
		})
	}

	// bottom
	if !(c.C.x == b.Height-1) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x + 1,
			y: c.C.y,
		})
	}

	// bottom left
	if c.C.y != 0 && c.C.x != b.Width-1 {
		neighbors = append(neighbors, coordinates{
			x: c.C.x + 1,
			y: c.C.y - 1,
		})
	}

	// left
	if !(c.C.x == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.C.x - 1,
			y: c.C.y,
		})
	}

	return neighbors
}

func (c *Cell) aliveNeighbors(b board) int {
	var count int
	for _, v := range c.neighbors(b) {
		if b.Cells[v.x][v.y].Alive {
			count++
		}
	}

	return count
}

func (c *Cell) die() {
	c.Alive = false
}

func (c *Cell) born() {
	c.Alive = true
}

type coordinates struct {
	x int
	y int
}

func (b *board) render() string {
	line := strings.Builder{}
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if b.Cells[i][j].Alive {
				line.Write([]byte(fmt.Sprintf("%c", aliveCell)))
			} else {
				line.Write([]byte(fmt.Sprintf("%c", deadCell)))
			}
		}
		line.Write([]byte("\n"))
	}
	return line.String()
}

func main() {
	end := make(chan os.Signal, 1)
	signal.Notify(end, os.Interrupt)

	b := newBoard(32, 32)
	b.Clock = 300 * time.Millisecond
	clock := time.Tick(b.Clock)

	// iterate the board and initialize a random state
	bg := newBoolgen()
	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {
			b.Cells[i][j].Alive = bg.Bool()
		}
	}

	// rectangle - oscillator
	//b.Cells[8][7].Alive = true
	//b.Cells[8][8].Alive = true
	//b.Cells[8][9].Alive = true

	// glider
	//b.Cells[3][8].Alive = true
	//b.Cells[4][9].Alive = true
	//b.Cells[5][9].Alive = true
	//b.Cells[5][8].Alive = true
	//b.Cells[5][7].Alive = true

	// still life
	//b.Cells[3][3].Alive = true
	//b.Cells[5][3].Alive = true
	//b.Cells[2][4].Alive = true
	//b.Cells[2][5].Alive = true
	//b.Cells[4][4].Alive = true
	//b.Cells[5][4].Alive = true
	//b.Cells[3][6].Alive = true
	//b.Cells[4][6].Alive = true
	//b.Cells[5][6].Alive = true
	//b.Cells[5][7].Alive = true

	go func() {
		for {
			select {
			case <-clock:
				fmt.Print("\033[H\033[2J")
				fmt.Println(b.render())
				b.runCycle()
			}
		}
	}()
	<-end
}

type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func newBoolgen() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}
