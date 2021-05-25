package main

import (
	"fmt"
	"strings"
	"time"
)

type board struct {
	width  int
	height int
	cells  [][]cell
	clock  time.Duration // how often a cycle lasts
}

func (b *board) next() {
	cells := make([][]cell, b.height)
	for i := 0; i < b.height; i++ {
		cells[i] = make([]cell, b.width)
		for j := 0; j < b.width; j++ {
			cells[i][j] = cell{
				alive: b.cells[i][j].alive,
				coordinates: coordinates{
					x: b.cells[i][j].coordinates.x,
					y: b.cells[i][j].coordinates.y,
				},
			}
		}
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			aliveNeighbors := b.cells[i][j].aliveNeighbors(*b)
			if b.cells[i][j].alive {
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

	b.cells = cells
}

func newBoard(width, height int) board {
	b := board{
		width:  width,
		height: height,
		cells:  make([][]cell, height),
	}

	for i := 0; i < height; i++ {
		b.cells[i] = make([]cell, width)
		for j := 0; j < width; j++ {
			b.cells[i][j].alive = false
			b.cells[i][j].coordinates.x = i
			b.cells[i][j].coordinates.y = j
		}
	}
	return b
}

func (b *board) render() string {
	line := strings.Builder{}
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.cells[i][j].alive {
				line.Write([]byte(fmt.Sprintf("%c", aliveCell)))
			} else {
				line.Write([]byte(fmt.Sprintf("%c", deadCell)))
			}
		}
		line.Write([]byte("\n"))
	}
	return line.String()
}
