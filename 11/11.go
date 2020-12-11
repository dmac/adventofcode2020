package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	input := "input.txt"

	grid, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}
	grid = simUntilSame(grid, adjacent, 4)
	part1 := countOccupied(grid)

	grid, err = readInput(input)
	if err != nil {
		log.Fatal(err)
	}
	grid = simUntilSame(grid, adjacentDiag, 5)
	part2 := countOccupied(grid)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func countOccupied(grid [][]byte) int {
	cnt := 0
	for _, row := range grid {
		for _, c := range row {
			if c == '#' {
				cnt++
			}
		}
	}
	return cnt
}

type adjFunc func([][]byte, int, int) int

func simUntilSame(grid [][]byte, afn adjFunc, numToEmpty int) [][]byte {
	for {
		ngrid := simOne(grid, afn, numToEmpty)
		if reflect.DeepEqual(grid, ngrid) {
			return grid
		}
		grid = ngrid
	}
}

func simOne(grid [][]byte, afn adjFunc, numToEmpty int) [][]byte {
	var next [][]byte
	for r := range grid {
		row := grid[r]
		nrow := make([]byte, len(row))
		for c := range row {
			switch row[c] {
			case '.':
				nrow[c] = '.'
			case 'L':
				adj := afn(grid, r, c)
				if adj == 0 {
					nrow[c] = '#'
				} else {
					nrow[c] = 'L'
				}
			case '#':
				adj := afn(grid, r, c)
				if adj >= numToEmpty {
					nrow[c] = 'L'
				} else {
					nrow[c] = '#'
				}
			default:
				panic("unexpected char")
			}
		}
		next = append(next, nrow)
	}
	return next
}

func adjacent(grid [][]byte, row, col int) int {
	cnt := 0
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || (r == row && c == col) {
				continue
			}
			if grid[r][c] == '#' {
				cnt++
			}
		}
	}
	return cnt
}

func adjacentDiag(grid [][]byte, row, col int) int {
	cnt := 0
	for r, c := row-1, col-1; r >= 0 && c >= 0; r, c = r-1, c-1 {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row-1, col; r >= 0; r-- {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row-1, col+1; r >= 0 && c < len(grid[0]); r, c = r-1, c+1 {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row, col-1; c >= 0; c-- {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row, col+1; c < len(grid[0]); c++ {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row+1, col-1; r < len(grid) && c >= 0; r, c = r+1, c-1 {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row+1, col; r < len(grid); r++ {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	for r, c := row+1, col+1; r < len(grid) && c < len(grid[0]); r, c = r+1, c+1 {
		if grid[r][c] == 'L' {
			break
		}
		if grid[r][c] == '#' {
			cnt++
			break
		}
	}
	return cnt
}

func readInput(filename string) ([][]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var grid [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return grid, nil
}
