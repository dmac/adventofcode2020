package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	seats, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	min := 99999999
	max := 0
	ids := make(map[int]struct{})
	for _, s := range seats {
		if s.id() < min {
			min = s.id()
		}
		if s.id() > max {
			max = s.id()
		}
		ids[s.id()] = struct{}{}
	}
	var mine int
	for i := min; i <= max; i++ {
		if _, ok := ids[i]; !ok {
			mine = i
			break
		}
	}
	fmt.Printf("Part 1: %d\n", max)
	fmt.Printf("Part 2: %d\n", mine)
}

type seat struct {
	row int
	col int
}

func (s seat) id() int { return s.row*8 + s.col }

func readInput(filename string) ([]seat, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var seats []seat
	for scanner.Scan() {
		line := scanner.Text()
		minr := 0
		maxr := 128
		minc := 0
		maxc := 8
		for _, c := range line[:7] {
			if c == 'F' {
				maxr -= (maxr - minr) / 2
			} else if c == 'B' {
				minr += (maxr - minr) / 2
			}
		}
		for _, c := range line[7:] {
			if c == 'L' {
				maxc -= (maxc - minc) / 2
			} else if c == 'R' {
				minc += (maxc - minc) / 2
			}
		}
		seats = append(seats, seat{minr, minc})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}
