package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputPrefix = "input"

func main() {
	keys, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	cardLoops := findLoopSize(keys[0])

	part1 := 1
	for i := 0; i < cardLoops; i++ {
		part1 = doLoop(keys[1], part1)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", 0)
}

func findLoopSize(key int) int {
	v := 1
	for i := 1; ; i++ {
		v = doLoop(7, v)
		if v == key {
			return i
		}
	}
	return 0
}

func doLoop(subject, v int) int {
	v *= subject
	return v % 20201227
}

func readInput(filename string) ([]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var keys []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		keys = append(keys, n)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}
