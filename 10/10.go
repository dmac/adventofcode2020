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
	nums = append(nums, 0)
	sort.Ints(nums)
	nums = append(nums, nums[len(nums)-1]+3)
	diffs := differences(nums)
	part1 := diffs[1] * diffs[3]
	part2 := arrangements(nums)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func differences(nums []int) (diffs [4]int) {
	curr := 0
	for _, n := range nums {
		if n-curr > 3 {
			panic("difference larger than 3")
		}
		diffs[n-curr]++
		curr = n
	}
	return
}

func arrangements(nums []int) int {
	var next [][]int
	for i := 0; i < len(nums); i++ {
		var n []int
		for j := i + 1; j < len(nums) && nums[j]-nums[i] <= 3; j++ {
			n = append(n, j)
		}
		next = append(next, n)
	}
	ans := make([]int, len(nums))
	ans[len(ans)-1] = 1
	for i := len(ans) - 2; i >= 0; i-- {
		for _, j := range next[i] {
			ans[i] += ans[j]
		}
	}
	return ans[0]
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
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		nums = append(nums, value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nums, nil
}
