package main

type cell struct {
	alive bool
	coordinates
}

type coordinates struct {
	x int
	y int
}

// neighbors returns all the neighbors of a cell for on a board b
func (c *cell) neighbors(b board) []coordinates {
	neighbors := make([]coordinates, 0)
	// top left
	if !(c.y == 0) && !(c.x == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.x - 1,
			y: c.y - 1,
		})
	}

	// top neighbor
	if !(c.y == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.x,
			y: c.y - 1,
		})
	}

	// top right
	if !(c.x == 0) && !(c.y == b.width-1) {
		neighbors = append(neighbors, coordinates{
			x: c.x - 1,
			y: c.y + 1,
		})
	}

	// right
	if !(c.y == b.width-1) {
		neighbors = append(neighbors, coordinates{
			x: c.x,
			y: c.y + 1,
		})
	}

	// bottom right
	if !(c.y == b.width-1) && !(c.x == b.height-1) {
		neighbors = append(neighbors, coordinates{
			x: c.x + 1,
			y: c.y + 1,
		})
	}

	// bottom
	if !(c.x == b.height-1) {
		neighbors = append(neighbors, coordinates{
			x: c.x + 1,
			y: c.y,
		})
	}

	// bottom left
	if c.y != 0 && c.x != b.width-1 {
		neighbors = append(neighbors, coordinates{
			x: c.x + 1,
			y: c.y - 1,
		})
	}

	// left
	if !(c.x == 0) {
		neighbors = append(neighbors, coordinates{
			x: c.x - 1,
			y: c.y,
		})
	}

	return neighbors
}

func (c *cell) aliveNeighbors(b board) int {
	var count int
	for _, v := range c.neighbors(b) {
		if b.cells[v.x][v.y].alive {
			count++
		}
	}

	return count
}

func (c *cell) die() {
	c.alive = false
}

func (c *cell) born() {
	c.alive = true
}
