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
	in, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	valid, part1 := findValidTickets(in)
	mapping := determineFieldOrder(in, valid)
	part2 := 1
	for i, idx := range mapping {
		f := in.fields[idx]
		if strings.HasPrefix(f.name, "departure") {
			part2 *= in.myTicket[i]
		}
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func determineFieldOrder(in *input, valid []ticket) []int {
	if len(in.fields) != len(in.myTicket) {
		panic("len fields != len ticket")
	}

	possibles := make([]map[int]struct{}, len(in.myTicket))
	for i := range possibles {
		possibles[i] = make(map[int]struct{})
		for j := 0; j < len(in.fields); j++ {
			possibles[i][j] = struct{}{}
		}
	}

	tickets := append([]ticket{in.myTicket}, valid...)
	for _, t := range tickets {
		for i, v := range t {
			for j, f := range in.fields {
				valid := false
				for _, r := range f.ranges {
					if v >= r[0] && v <= r[1] {
						valid = true
					}
				}
				if !valid {
					delete(possibles[i], j)
				}
			}
		}
	}

	locked := make(map[int]struct{})
	for i := 0; i < len(possibles); i++ {
		p := possibles[i]
		if len(p) == 0 {
			panic("bad state")
		}
		if len(p) > 1 {
			continue
		}
		if _, ok := locked[i]; ok {
			continue
		}
		var v int
		for vv := range p {
			v = vv
		}
		for j := 0; j < len(possibles); j++ {
			if j == i {
				continue
			}
			delete(possibles[j], v)
		}
		locked[i] = struct{}{}
		i = -1
	}

	var mapping []int
	for _, p := range possibles {
		var v int
		for vv := range p {
			v = vv
		}
		mapping = append(mapping, v)
	}
	return mapping
}

func findValidTickets(in *input) (valid []ticket, sum int) {
	var invalid []int
	for _, t := range in.nearbyTickets {
		validT := true
		for _, v := range t {
			validV := false
			for _, f := range in.fields {
				for _, r := range f.ranges {
					if v >= r[0] && v <= r[1] {
						validV = true
						goto checkValid
					}
				}
			}
		checkValid:
			if !validV {
				validT = false
				invalid = append(invalid, v)
			}
		}
		if validT {
			valid = append(valid, t)
		}
	}
	for _, v := range invalid {
		sum += v
	}
	return valid, sum
}

type input struct {
	fields        []field
	myTicket      ticket
	nearbyTickets []ticket
}

type field struct {
	name   string
	ranges [][2]int
}

type ticket []int

func readInput(filename string) (*input, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var in input
	scannedFields := false
	scannedMyTicket := false
	for scanner.Scan() {
		line := scanner.Text()
		if !scannedFields {
			if line == "" {
				scannedFields = true
				continue
			}
			parts := strings.Split(line, ": ")
			if len(parts) != 2 {
				panic("len parts != 2")
			}
			field := field{name: parts[0]}
			ranges := strings.Split(parts[1], " or ")
			for _, r := range ranges {
				lohi := strings.Split(r, "-")
				lo, err := strconv.Atoi(lohi[0])
				if err != nil {
					return nil, err
				}
				hi, err := strconv.Atoi(lohi[1])
				if err != nil {
					return nil, err
				}
				field.ranges = append(field.ranges, [2]int{lo, hi})
			}
			in.fields = append(in.fields, field)
		}
		if strings.Contains(line, ":") {
			continue
		}
		if !scannedMyTicket {
			if line == "" {
				scannedMyTicket = true
				continue
			}
			parts := strings.Split(line, ",")
			for _, p := range parts {
				n, err := strconv.Atoi(p)
				if err != nil {
					return nil, err
				}
				in.myTicket = append(in.myTicket, n)
			}
			continue
		}
		var t ticket
		parts := strings.Split(line, ",")
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			t = append(t, n)
		}
		in.nearbyTickets = append(in.nearbyTickets, t)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &in, nil
}
