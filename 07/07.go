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
	rules, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}
	bagsThatCanContain := make(map[string][]string)
	bagsContainedBy := make(map[string][]bagNum)
	for _, r := range rules {
		bagsContainedBy[r.bag] = r.contains
		for _, c := range r.contains {
			bagsThatCanContain[c.bag] = append(bagsThatCanContain[c.bag], r.bag)
		}
	}

	part1 := make(map[string]struct{})
	seen := make(map[string]struct{})
	queue := []string{"shiny gold"}
	for len(queue) > 0 {
		bag := queue[0]
		queue = queue[1:]
		if _, ok := seen[bag]; ok {
			continue
		}
		seen[bag] = struct{}{}
		bags := bagsThatCanContain[bag]
		for _, b := range bags {
			part1[b] = struct{}{}
			queue = append(queue, b)
		}
	}

	part2 := 0
	queue = []string{"shiny gold"}
	for len(queue) > 0 {
		bag := queue[0]
		queue = queue[1:]
		for _, contained := range bagsContainedBy[bag] {
			part2 += contained.num
			for i := 0; i < contained.num; i++ {
				queue = append(queue, contained.bag)
			}
		}
	}

	fmt.Printf("Part 1: %d\n", len(part1))
	fmt.Printf("Part 2: %d\n", part2)
}

type rule struct {
	bag      string
	contains []bagNum
}

type bagNum struct {
	num int
	bag string
}

func readInput(filename string) ([]rule, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rules []rule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " contain ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expect len(parts) == 2")
		}
		bag := strings.TrimSuffix(parts[0], " bags")
		parts = strings.Split(parts[1], ", ")
		r := rule{bag: bag}
		for _, part := range parts {
			i := strings.Index(part, " ")
			if i < 0 {
				return nil, fmt.Errorf("expect i >= 0")
			}
			var n int
			if part[:i] == "no" {
				continue
			}
			n, err := strconv.Atoi(part[:i])
			if err != nil {
				return nil, err
			}
			part = part[i+1:]
			j := strings.Index(part, " bag")
			if j < 0 {
				return nil, fmt.Errorf("expect i >= 0")
			}
			nb := bagNum{
				bag: part[:j],
				num: n,
			}
			r.contains = append(r.contains, nb)
		}
		rules = append(rules, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}
