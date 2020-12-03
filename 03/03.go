package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	lines, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1: %d\n", treesEncountered(lines, 3, 1))

	ans := treesEncountered(lines, 1, 1)
	ans *= treesEncountered(lines, 3, 1)
	ans *= treesEncountered(lines, 5, 1)
	ans *= treesEncountered(lines, 7, 1)
	ans *= treesEncountered(lines, 1, 2)
	fmt.Printf("Part 2: %d\n", ans)
}

func treesEncountered(lines []string, right, down int) int {
	row := 0
	col := 0
	count := 0
	for row < len(lines) {
		if lines[row][col] == '#' {
			count++
		}
		row += down
		col += right
		if col >= len(lines[0]) {
			col %= len(lines[0])
		}
	}
	return count
}

func readInput(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
