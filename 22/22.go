package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const inputPrefix = "input"

func main() {
	log.SetFlags(0)
	p1, p2, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	p1, p2 = playGame(p1, p2)

	part1 := 0
	if len(p1) == 0 {
		part1 = scoreDeck(p2)
	} else {
		part1 = scoreDeck(p1)
	}

	p1, p2, err = readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	part2 := 0
	var winner int
	p1, p2, winner = playGameRecursive(p1, p2)
	if winner == 1 {
		part2 = scoreDeck(p1)
	} else {
		part2 = scoreDeck(p2)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func playGameRecursive(p1, p2 []int) ([]int, []int, int) {
	states := make(map[string]struct{})
	for len(p1) > 0 && len(p2) > 0 {
		hash := deckHash(p1) + "|" + deckHash(p2)
		if _, ok := states[hash]; ok {
			return p1, p2, 1
		}
		states[hash] = struct{}{}
		var c1, c2 int
		c1, p1 = p1[0], p1[1:]
		c2, p2 = p2[0], p2[1:]
		if c1 <= len(p1) && c2 <= len(p2) {
			p3 := make([]int, c1)
			p4 := make([]int, c2)
			copy(p3, p1)
			copy(p4, p2)
			_, _, winner := playGameRecursive(p3, p4)
			switch winner {
			case 1:
				p1 = append(p1, c1, c2)
			case 2:
				p2 = append(p2, c2, c1)
			default:
				panic("unknown winner of subgame")
			}
			continue
		}
		switch {
		case c1 > c2:
			p1 = append(p1, c1, c2)
		case c2 > c1:
			p2 = append(p2, c2, c1)
		default:
			panic("cards equal")
		}
	}
	var winner int
	if len(p1) == 0 {
		winner = 2
	} else {
		winner = 1
	}
	return p1, p2, winner
}

func deckHash(deck []int) string {
	var sb strings.Builder
	for _, n := range deck {
		sb.WriteString(fmt.Sprintf("%d,", n))
	}
	return sb.String()
}

func scoreDeck(deck []int) int {
	score := 0
	for i := 1; i <= len(deck); i++ {
		score += i * deck[len(deck)-i]
	}
	return score
}

func playGame(p1, p2 []int) ([]int, []int) {
	for len(p1) > 0 && len(p2) > 0 {
		var c1, c2 int
		c1, p1 = p1[0], p1[1:]
		c2, p2 = p2[0], p2[1:]
		switch {
		case c1 > c2:
			p1 = append(p1, c1, c2)
		case c2 > c1:
			p2 = append(p2, c2, c1)
		default:
			panic("cards equal")
		}
	}
	return p1, p2
}

func readInput(filename string) (p1, p2 []int, _ error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	p1Done := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			p1Done = true
			continue
		}
		if strings.Contains(line, "Player") {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, nil, err
		}
		if p1Done {
			p2 = append(p2, n)
		} else {
			p1 = append(p1, n)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return p1, p2, nil
}
