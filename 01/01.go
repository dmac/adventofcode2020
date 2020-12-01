package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	nums, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
l1:
	for i, a := range nums {
		for _, b := range nums[i+1:] {
			if a+b == 2020 {
				fmt.Printf(
					"%d + %d = 2020; %d * %d = %d\n",
					a, b, a, b, a*b,
				)
				break l1
			}
		}
	}
l2:
	for i, a := range nums {
		for j, b := range nums[i+1:] {
			for _, c := range nums[j+1:] {
				if a+b+c == 2020 {
					fmt.Printf(
						"%d + %d + %d = 2020; %d * %d * %d = %d\n",
						a, b, c, a, b, c, a*b*c,
					)
					break l2
				}
			}
		}
	}
}

func readInput(filename string) ([]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var nums []int
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
