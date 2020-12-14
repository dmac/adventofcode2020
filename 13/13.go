package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	earliest, ids, err := readInput("small.txt")
	if err != nil {
		log.Fatal(err)
	}

	part1 := findFirst(earliest, ids)
	// part2 := findSecond(ids)
	part2 := findThird(ids)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func findFirst(earliest int, ids []int) int {
	times := make([]int, len(ids))
	for i, id := range ids {
		if id == 0 {
			continue
		}
		for j := earliest - (earliest % id); ; j += id {
			if j >= earliest {
				times[i] = j
				break
			}
		}
	}
	minTime := 999999999
	minID := 999999999
	for i, t := range times {
		if t == 0 {
			continue
		}
		if t < minTime {
			minTime = t
			minID = ids[i]
		}
	}
	return minID * (minTime - earliest)
}

type target struct {
	id  int
	rem int
}

func findSecond(ids []int) int {
	var targets []target
	for i, id := range ids {
		if i == 0 || id == 0 {
			continue
		}
		targets = append(targets, target{id: id, rem: id - i})
	}
outer:
	// for i := 0; ; i += ids[0] {
	for i := 100000000000000; ; i += ids[0] {
		if i == 1068781 {
			fmt.Println(i)
		}
		for _, tar := range targets {
			if i%tar.id != tar.rem {
				continue outer
			}
		}
		return i
	}
}

func findThird(ids []int) int {
	var targets []target
	for i, id := range ids {
		if id == 0 {
			continue
		}
		if i == 0 {
			targets = append(targets, target{id: id, rem: 0})
		} else {
			targets = append(targets, target{id: id, rem: id - i})
		}
	}
	fmt.Println(targets)
	return 0
}

func readInput(filename string) (int, []int, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, nil, err
	}
	lines := bytes.Split(b, []byte("\n"))
	earliest, err := strconv.Atoi(string(lines[0]))
	if err != nil {
		return 0, nil, err
	}
	parts := bytes.Split(lines[1], []byte(","))
	var ids []int
	for _, part := range parts {
		if string(part) == "x" {
			ids = append(ids, 0)
			continue
		}
		id, err := strconv.Atoi(string(part))
		if err != nil {
			return 0, nil, err
		}
		ids = append(ids, id)
	}
	return earliest, ids, nil
}
