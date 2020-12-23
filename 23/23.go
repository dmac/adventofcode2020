package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const inputPrefix = "input"

func main() {
	log.SetFlags(0)

	cups, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	part1 := doPart1LL(cups)
	fmt.Printf("Part 1: %s\n", part1)

	cups, err = readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := 10; i <= 1_000_000; i++ {
		cups = append(cups, i)
	}
	part2 := doPart2LL(cups)
	fmt.Printf("Part 2: %d\n", part2)
}

func doPart2(cups []int) int {
	cups = simulate(cups, 10_000_000)
	for i, c := range cups {
		if c == 1 {
			j := (i + 1) % len(cups)
			k := (i + 2) % len(cups)
			return cups[j] * cups[k]
		}
	}
	return -1
}

func doPart2LL(cups []int) int {
	head := simulateLL(cups, 10_000_000)
	for ; ; head = head.next {
		if head.label == 1 {
			a := head.next.label
			b := head.next.next.label
			return a * b
		}
	}
	return -1
}

func doPart1(cups []int) string {
	cups = simulate(cups, 100)
	var nums []int
	for i, c := range cups {
		if c == 1 {
			nums = append(nums, cups[i+1:]...)
			nums = append(nums, cups[:i]...)
			break
		}
	}
	var sb strings.Builder
	for _, c := range nums {
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	return sb.String()
}

func doPart1LL(cups []int) string {
	head := simulateLL(cups, 100)
	var sb strings.Builder
	one := false
	for ; !(one && head.label == 1); head = head.next {
		if one {
			sb.WriteString(fmt.Sprintf("%d", head.label))
		}
		if head.label == 1 {
			one = true
		}
	}
	return sb.String()
}

func simulate(cups []int, moves int) []int {
	min := math.MaxInt64
	max := 0
	for _, c := range cups {
		if c < min {
			min = c
		}
		if c > max {
			max = c
		}
	}
	idx := 0
	var outs [3]int
	for i := 0; i < moves; i++ {
		curr := cups[idx]
		out := (idx + 1) % max
		for j := 0; j < 3; j++ {
			outs[j] = cups[(out+j)%max]
		}
		dest := curr - 1
		for j := out; j != (out+3)%max; j = (j + 1) % max {
			if dest < min {
				dest = max
			}
			if dest == cups[j] {
				dest--
				j = out - 1
			}
		}
		from := out - 1
		to := (out + 2) % max
		for {
			if from < 0 {
				from = len(cups) - 1
			}
			if to < 0 {
				to = len(cups) - 1
			}
			if cups[from] == dest {
				break
			}
			cups[to] = cups[from]
			from--
			to--
		}
		for j := 2; j >= 0; j-- {
			if to < 0 {
				to = len(cups) - 1
			}
			cups[to] = outs[j]
			to--
		}
		for j, c := range cups {
			if c == curr {
				idx = (j + 1) % len(cups)
				break
			}
		}
	}
	return cups
}

func simulateLL(cups []int, moves int) *cup {
	head := makeCups(cups)
	min := math.MaxInt64
	max := 0
	index := make(map[int]*cup)
	for c := head.next; ; c = c.next {
		index[c.label] = c
		if c.label < min {
			min = c.label
		}
		if c.label > max {
			max = c.label
		}
		if c == head {
			break
		}
	}
	curr := head
	for i := 0; i < moves; i++ {
		outs := curr.next
		for j := 0; j < 3; j++ {
			curr.next = curr.next.next
		}
		dest := curr.label - 1
		c := outs
		for j := 0; j < 3; j++ {
			if dest < min {
				dest = max
			}
			if dest == c.label {
				dest--
				c = outs
				j = -1
			} else {
				c = c.next
			}
		}
		c = index[dest]
		tail := c.next
		c.next = outs
		for j := 0; j < 3; j++ {
			c = c.next
		}
		c.next = tail
		curr = curr.next
	}
	return curr
}

func makeCups(cups []int) *cup {
	backing := make([]cup, len(cups))
	head := &backing[0]
	head.label = cups[0]
	curr := head
	for i := 1; i < len(cups); i++ {
		c := &backing[i]
		c.label = cups[i]
		curr.next = c
		curr = c
	}
	curr.next = head
	return head
}

type cup struct {
	label int
	next  *cup
}

func readInput(filename string) ([]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cups []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for _, r := range line {
			n, err := strconv.Atoi(string([]rune{r}))
			if err != nil {
				return nil, err
			}
			cups = append(cups, n)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cups, nil
}
