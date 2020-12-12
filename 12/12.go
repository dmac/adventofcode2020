package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	input := "input.txt"

	dirs, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}

	sh := &ship{orient: E}
	sh.navigateShip(dirs)
	part1 := abs(sh.x) + abs(sh.y)

	sh = &ship{
		orient: E,
		wx:     10,
		wy:     1,
	}
	sh.navigateWaypoint(dirs)
	part2 := abs(sh.x) + abs(sh.y)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type orientation int

const (
	N = iota
	E
	S
	W
)

func (o orientation) String() string {
	switch o {
	case N:
		return "N"
	case E:
		return "E"
	case S:
		return "S"
	case W:
		return "W"
	default:
		return "?"
	}
}

type ship struct {
	x, y   int
	wx, wy int
	orient orientation
}

func (s *ship) String() string {
	return fmt.Sprintf("x=%d y=%d o=%s", s.x, s.y, s.orient)
}

func (s *ship) navigateShip(dirs []direction) {
	for _, dir := range dirs {
		switch dir.action {
		case 'N':
			s.y += dir.value
		case 'S':
			s.y -= dir.value
		case 'E':
			s.x += dir.value
		case 'W':
			s.x -= dir.value
		case 'L':
			turns := dir.value / 90
			o := int(s.orient)
			for i := 0; i < turns; i++ {
				o--
				if o < 0 {
					o = 3
				}
			}
			s.orient = orientation(o)
		case 'R':
			turns := dir.value / 90
			o := int(s.orient)
			for i := 0; i < turns; i++ {
				o++
				if o >= 3 {
					o = 0
				}
			}
			s.orient = orientation(abs(int(s.orient)+turns) % 4)
		case 'F':
			switch s.orient {
			case N:
				s.y += dir.value
			case E:
				s.x += dir.value
			case S:
				s.y -= dir.value
			case W:
				s.x -= dir.value
			default:
				panic("bad direction")
			}
		}
	}
}

func (s *ship) navigateWaypoint(dirs []direction) {
	for _, dir := range dirs {
		switch dir.action {
		case 'N':
			s.wy += dir.value
		case 'S':
			s.wy -= dir.value
		case 'E':
			s.wx += dir.value
		case 'W':
			s.wx -= dir.value
		case 'L':
			turns := dir.value / 90
			for i := 0; i < turns; i++ {
				s.wx, s.wy = -s.wy, s.wx
			}
		case 'R':
			turns := dir.value / 90
			for i := 0; i < turns; i++ {
				s.wx, s.wy = s.wy, -s.wx
			}
		case 'F':
			for i := 0; i < dir.value; i++ {
				s.x += s.wx
				s.y += s.wy
			}
		}
	}
}

type direction struct {
	action byte
	value  int
}

func readInput(filename string) ([]direction, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dirs []direction
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		v, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		dir := direction{
			action: line[0],
			value:  v,
		}
		dirs = append(dirs, dir)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return dirs, nil
}
