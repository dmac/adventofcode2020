package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	input := "input.txt"
	groups, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}
	any := 0
	every := 0
	for _, group := range groups {
		m := make(map[rune]int)
		for _, person := range group {
			for _, r := range person {
				m[r]++
			}
		}
		any += len(m)
		for _, n := range m {
			if n == len(group) {
				every++
			}
		}
	}
	fmt.Printf("Part 1: %d\n", any)
	fmt.Printf("Part 2: %d\n", every)
}

func readInput(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var groups [][]string
	var group []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			group = append(group, line)
			continue
		}
		groups = append(groups, group)
		group = nil
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	groups = append(groups, group)
	return groups, nil
}
