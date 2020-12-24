package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

const inputPrefix = "input"

func main() {
	dirs, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	f := &floor{grid: make(map[[2]float64]bool)}
	part1 := f.flipTiles(dirs)
	part2 := f.doArt()

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

type floor struct {
	black int
	grid  map[[2]float64]bool
}

func (f *floor) doArt() int {
	for i := 0; i < 100; i++ {
		next := make(map[[2]float64]bool)
		minr, minc := math.MaxFloat64, math.MaxFloat64
		maxr, maxc := float64(0), float64(0)
		for tile := range f.grid {
			if tile[0] < minr {
				minr = tile[0]
			}
			if tile[0] > maxr {
				maxr = tile[0]
			}
			if tile[1] < minc {
				minc = tile[1]
			}
			if tile[1] > maxc {
				maxc = tile[1]
			}
		}
		for r := minr - 1; r <= maxr+1; r++ {
			for c := minc - 1; c <= maxc+1; c += 0.5 {
				n := 0
				tile := [2]float64{r, c}
				for _, d := range [][2]float64{
					{-1, -0.5},
					{-1, 0.5},
					{0, 1},
					{1, 0.5},
					{1, -0.5},
					{0, -1},
				} {
					b := f.grid[[2]float64{tile[0] + d[0], tile[1] + d[1]}]
					if b {
						n++
					}
				}
				black := f.grid[tile]
				if black && (n == 0 || n > 2) {
					next[tile] = false
					f.black--
				} else if !black && n == 2 {
					next[tile] = true
					f.black++
				} else {
					next[tile] = black
				}
			}
		}
		f.grid = next
	}
	return f.black
}

func (f *floor) flipTiles(dirs [][]string) int {
	for _, dir := range dirs {
		r := float64(0)
		c := float64(0)
		for _, d := range dir {
			switch d {
			case "e":
				c++
			case "w":
				c--
			case "ne":
				r--
				c += 0.5
			case "nw":
				r--
				c -= 0.5
			case "se":
				r++
				c += 0.5
			case "sw":
				r++
				c -= 0.5
			}
		}
		tile := [2]float64{r, c}
		black := f.grid[tile]
		black = !black
		f.grid[tile] = black
		if black {
			f.black++
		} else {
			f.black--
		}
	}
	return f.black
}

func readInput(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dirs [][]string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var dir []string
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		for i := 0; i < len(line); i++ {
			switch line[i] {
			case 'e', 'w':
				dir = append(dir, string(line[i]))
			case 'n', 's':
				dir = append(dir, string([]byte{line[i], line[i+1]}))
				i++
			}
		}
		dirs = append(dirs, dir)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return dirs, nil
}
