package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
)

func main() {
	earliest, ids, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	part1 := findFirst(earliest, ids)
	part2 := findSecond(ids)

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
	minTime := math.MaxInt64
	minID := math.MaxInt64
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

func findSecond(ids []int) int {
	sofar := ids[0]
	var t int
	for i := 1; i < len(ids); i++ {
		id := ids[i]
		if id == 0 {
			continue
		}
		for ; ; t += sofar {
			if (t+i)%id == 0 {
				sofar *= id
				break
			}
		}
	}
	return t
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
