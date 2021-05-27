package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

var (
	aliveCell = '\u2B1B'
	deadCell  = '\u2B1C'
)

func main() {
	end := make(chan os.Signal, 1)
	signal.Notify(end, os.Interrupt)

	b := newBoard(25, 25)
	b.clock = 200 * time.Millisecond
	clock := time.Tick(b.clock)

	// iterate the board and initialize a random state
	bg := newBoolgen()
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			b.cells[i][j].alive = bg.Bool()
		}
	}

	// rectangle - oscillator
	//b.cells[8][7].alive = true
	//b.cells[8][8].alive = true
	//b.cells[8][9].alive = true

	// glider
	b.cells[3][8].alive = true
	b.cells[4][9].alive = true
	b.cells[5][9].alive = true
	b.cells[5][8].alive = true
	b.cells[5][7].alive = true

	go func() {
		for {
			select {
			case <-clock:
				fmt.Print("\033[H\033[2J")
				fmt.Println(b.render())
				b.next()
			}
		}
	}()
	<-end
}
