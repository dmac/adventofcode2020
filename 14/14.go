package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	insts, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	mem := make(map[int]int)
	runV1(mem, insts)
	part1 := 0
	for _, v := range mem {
		part1 += v
	}

	mem = make(map[int]int)
	runV2(mem, insts)
	part2 := 0
	for _, v := range mem {
		part2 += v
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func runV1(mem map[int]int, insts []instruction) {
	var mask mask
	for _, inst := range insts {
		if inst.isMask {
			mask = inst.mask
			continue
		}
		mem[inst.addr] = applyMaskV1(mask, inst.val)
	}
}

func runV2(mem map[int]int, insts []instruction) {
	var mask mask
	for _, inst := range insts {
		if inst.isMask {
			mask = inst.mask
			continue
		}
		applyMaskV2(mem, mask, uint64(inst.addr), uint64(inst.val))
	}
}

func applyMaskV1(mask mask, val int) int {
	v := uint64(val)
	v |= mask.set
	v &= ^mask.unset
	return int(v)
}

func applyMaskV2(mem map[int]int, mask mask, addr, val uint64) {
	if len(mask.xs) == 0 {
		return
	}
	x := mask.xs[0]
	mask.xs = mask.xs[1:]
	addr |= mask.set
	addr &^= 1 << x
	mem[int(addr)] = int(val)
	applyMaskV2(mem, mask, addr, val)
	addr |= 1 << x
	mem[int(addr)] = int(val)
	applyMaskV2(mem, mask, addr, val)
}

type instruction struct {
	addr int
	val  int

	mask   mask
	isMask bool
}

type mask struct {
	set   uint64
	unset uint64
	xs    []int
}

func parseMask(m string) mask {
	var mask mask
	for i := 0; i < len(m); i++ {
		b := m[len(m)-i-1]
		if b == '1' {
			mask.set |= 1 << i
		}
		if b == '0' {
			mask.unset |= 1 << i
		}
		if b == 'X' {
			mask.xs = append(mask.xs, i)
		}
	}
	return mask
}

func parseMemSet(line string) (int, int) {
	var addr, val int
	if _, err := fmt.Sscanf(line, "mem[%d] = %d", &addr, &val); err != nil {
		panic(err)
	}
	return addr, val
}

func readInput(filename string) ([]instruction, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var insts []instruction
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var inst instruction
		line := scanner.Text()
		if strings.HasPrefix(line, "mask = ") {
			inst.mask = parseMask(line[len("mask = "):])
			inst.isMask = true
		} else {
			inst.addr, inst.val = parseMemSet(line)
		}
		insts = append(insts, inst)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return insts, nil
}
