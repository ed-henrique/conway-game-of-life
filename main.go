package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TPF = 130
const GRID_SIZE = 50

type Grid [GRID_SIZE][GRID_SIZE]bool

func (g *Grid) Print() {
	for _, line := range g {
		for _, cell := range line {
			if cell {
				fmt.Print("\x1b[91m\x1b[101mX \x1b[39m\x1b[49m")
			} else {
				fmt.Print("  ")
			}
		}

		fmt.Println("")
	}

	fmt.Printf("\x1b[%dA", GRID_SIZE)
}

func NewGrid(oldGrid *Grid) *Grid {
	grid := &Grid{}

	if oldGrid == nil {
		for i := range GRID_SIZE {
			for j := range GRID_SIZE {
				if rand.Float64() < 0.3 {
					grid[i][j] = !grid[i][j]
				}
			}
		}
	} else {
		for i := range GRID_SIZE {
			for j := range GRID_SIZE {
				grid[i][j] = isCellAlive(oldGrid, i, j)
			}
		}
	}

	return grid
}

func isCellAlive(grid *Grid, x, y int) bool {
	aliveNeighbors := 0

	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if x == i && y == j || (i < 0 || i > GRID_SIZE-1 || j < 0 || j > GRID_SIZE-1) {
				continue
			}

			if grid[i][j] {
				aliveNeighbors++
			}
		}
	}

	if grid[x][y] {
		return aliveNeighbors == 2 || aliveNeighbors == 3
	}

	return aliveNeighbors == 3
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\x1b[?25l")

	g := NewGrid(nil)

	go func() {
		for {
			g.Print()
			g = NewGrid(g)
			time.Sleep(TPF * time.Millisecond)
		}
	}()

	<-sigs
	fmt.Println("\x1b[2A\x1b[?25h")
}
