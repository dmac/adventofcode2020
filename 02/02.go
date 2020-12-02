package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	valid := 0
	for _, l := range lines {
		n := 0
		for i := range l.password {
			if l.password[i] == l.letter {
				n++
			}
		}
		if n >= l.lo && n <= l.hi {
			valid++
		}
	}
	fmt.Printf("Part 1: %d\n", valid)

	valid = 0
	for _, l := range lines {
		if l.password[l.lo-1] == l.letter && l.password[l.hi-1] != l.letter {
			valid++
		} else if l.password[l.lo-1] != l.letter && l.password[l.hi-1] == l.letter {
			valid++
		}
	}
	fmt.Printf("Part 2: %d\n", valid)
}

func readInput(filename string) ([]line, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []line
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			panic("len(fields) != 3")
		}
		parts := strings.Split(fields[0], "-")
		if len(parts) != 2 {
			panic("len(parts) != 2")
		}
		var l line
		var err error
		l.lo, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		l.hi, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		l.letter = fields[1][0]
		l.password = fields[2]
		lines = append(lines, l)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

type line struct {
	lo       int
	hi       int
	letter   byte
	password string
}
