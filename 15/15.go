package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	nums, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	part1 := playGame(nums, 2020)
	part2 := playGame(nums, 30000000)
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func playGame(nums []int, limit int) int {
	m := make(map[int][]int)
	last := -1
	for turn := 1; ; turn++ {
		if turn == limit+1 {
			return last
		}
		if turn-1 < len(nums) {
			n := nums[turn-1]
			m[n] = append(m[n], turn)
			last = n
			continue
		}
		prevs := m[last]
		if len(prevs) == 1 {
			m[0] = append(m[0], turn)
			last = 0
			continue
		}
		t := prevs[len(prevs)-2]
		n := turn - t - 1
		m[n] = append(m[n], turn)
		last = n
	}
}

func readInput(filename string) ([]int, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(b), ",")
	var nums []int
	for _, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}
