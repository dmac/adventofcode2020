package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	input := "input.txt"
	nums, err := readInput(input)
	if err != nil {
		log.Fatal(err)
	}

	invalid := findInvalid(nums, 25)
	contiguous := findContiguous(nums, invalid)
	sort.Ints(contiguous)
	part2 := contiguous[0] + contiguous[len(contiguous)-1]

	fmt.Printf("Part 1: %d\n", invalid)
	fmt.Printf("Part 2: %d\n", part2)
}

func findContiguous(nums []int, invalid int) []int {
	for i := 0; i < len(nums); i++ {
		sum := nums[i]
		for j := i + 1; j < len(nums) && sum < invalid; j++ {
			sum += nums[j]
			if sum == invalid {
				return nums[i : j+1]
			}
		}
	}
	return nil
}

func findInvalid(nums []int, preamble int) int {
	for i := preamble; i < len(nums); i++ {
		n := nums[i]
		if !findPair(nums[i-preamble:i], n) {
			return n
		}
	}
	return 0
}

func findPair(nums []int, target int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return true
			}
		}
	}
	return false
}

func readInput(filename string) ([]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var nums []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nums, nil
}
