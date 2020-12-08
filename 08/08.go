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
	input := "input.txt"
	instructions, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}
	r := newRunner()
	r.run(instructions)
	part1 := r.acc

	part2 := tryToFix(instructions)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

type instruction struct {
	name  string
	value int
}

type runner struct {
	acc int
	pc  int

	seen map[int]struct{}
}

func tryToFix(instructions []instruction) int {
	for i, inst := range instructions {
		if inst.name == "acc" {
			continue
		}
		try := make([]instruction, len(instructions))
		copy(try, instructions[:i])
		if inst.name == "nop" {
			try[i].name = "jmp"
		} else {
			try[i].name = "nop"
		}
		try[i].value = inst.value
		copy(try[i+1:], instructions[i+1:])
		r := newRunner()
		if r.run(try) {
			return r.acc
		}
	}
	return 0
}

func newRunner() *runner {
	return &runner{
		seen: make(map[int]struct{}),
	}
}

func (r *runner) run(instructions []instruction) bool {
	for {
		if r.pc == len(instructions) {
			return true
		}
		if _, ok := r.seen[r.pc]; ok {
			return false
		}
		r.seen[r.pc] = struct{}{}
		inst := instructions[r.pc]
		switch inst.name {
		case "acc":
			r.acc += inst.value
			r.pc++
		case "nop":
			r.pc++
		case "jmp":
			r.pc += inst.value
		}
	}
}

func readInput(filename string) ([]instruction, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var instructions []instruction
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expect len(parts) == 2")
		}
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		inst := instruction{
			name:  parts[0],
			value: value,
		}
		instructions = append(instructions, inst)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return instructions, nil
}
