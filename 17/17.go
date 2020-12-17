package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

func main() {
	g, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		g.step3()
	}
	part1 := len(g.cubes)

	g, err = readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		g.step4()
	}
	part2 := len(g.cubes)

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

type position struct {
	x, y, z, w int
}

type game struct {
	cubes map[position]struct{}
}

func (g *game) String() string {
	var sb strings.Builder
	min := g.min()
	max := g.max()
	fmt.Println(min, max)
	for z := min.z; z <= max.z; z++ {
		for y := min.y; y <= max.y; y++ {
			for x := min.x; x <= max.x; x++ {
				pos := position{x: x, y: y, z: z}
				if _, ok := g.cubes[pos]; ok {
					sb.WriteString("#")
				} else {
					sb.WriteString(".")
				}
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *game) step3() {
	next := make(map[position]struct{})
	min := g.min()
	max := g.max()
	for z := min.z - 1; z <= max.z+1; z++ {
		for y := min.y - 1; y <= max.y+1; y++ {
			for x := min.x - 1; x <= max.x+1; x++ {
				pos := position{x: x, y: y, z: z}
				if _, ok := g.cubes[pos]; ok {
					n := g.activeNeighbors3(pos)
					if n == 2 || n == 3 {
						next[pos] = struct{}{}
					}
				}
				if g.activeNeighbors3(pos) == 3 {
					next[pos] = struct{}{}
				}
			}
		}
	}
	g.cubes = next
}

func (g *game) step4() {
	next := make(map[position]struct{})
	min := g.min()
	max := g.max()
	for w := min.w - 1; w <= max.w+1; w++ {
		for z := min.z - 1; z <= max.z+1; z++ {
			for y := min.y - 1; y <= max.y+1; y++ {
				for x := min.x - 1; x <= max.x+1; x++ {
					pos := position{x: x, y: y, z: z, w: w}
					if _, ok := g.cubes[pos]; ok {
						n := g.activeNeighbors4(pos)
						if n == 2 || n == 3 {
							next[pos] = struct{}{}
						}
					}
					if g.activeNeighbors4(pos) == 3 {
						next[pos] = struct{}{}
					}
				}
			}
		}
	}
	g.cubes = next
}

func (g *game) activeNeighbors3(pos position) int {
	active := 0
	for dz := -1; dz <= 1; dz++ {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				cube := position{
					x: pos.x + dx,
					y: pos.y + dy,
					z: pos.z + dz,
				}
				if _, ok := g.cubes[cube]; ok {
					active++
				}
			}
		}
	}
	return active
}

func (g *game) activeNeighbors4(pos position) int {
	active := 0
	for dw := -1; dw <= 1; dw++ {
		for dz := -1; dz <= 1; dz++ {
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 && dz == 0 && dw == 0 {
						continue
					}
					cube := position{
						x: pos.x + dx,
						y: pos.y + dy,
						z: pos.z + dz,
						w: pos.w + dw,
					}
					if _, ok := g.cubes[cube]; ok {
						active++
					}
				}
			}
		}
	}
	return active
}

func (g *game) min() position {
	minX := math.MaxInt64
	minY := math.MaxInt64
	minZ := math.MaxInt64
	minW := math.MaxInt64
	for pos := range g.cubes {
		if pos.x < minX {
			minX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.z < minZ {
			minZ = pos.z
		}
		if pos.w < minW {
			minW = pos.w
		}
	}
	return position{x: minX, y: minY, z: minZ, w: minW}
}

func (g *game) max() position {
	maxX := math.MinInt64
	maxY := math.MinInt64
	maxZ := math.MinInt64
	maxW := math.MinInt64
	for pos := range g.cubes {
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y > maxY {
			maxY = pos.y
		}
		if pos.z > maxZ {
			maxZ = pos.z
		}
		if pos.w > maxW {
			maxW = pos.w
		}
	}
	return position{x: maxX, y: maxY, z: maxZ, w: maxW}
}

func readInput(filename string) (*game, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	g := &game{cubes: make(map[position]struct{})}
	for y := 0; y < len(lines); y++ {
		line := lines[y]
		for x := 0; x < len(line); x++ {
			if line[x] == '#' {
				g.cubes[position{x: x, y: y}] = struct{}{}
			}
		}
	}
	return g, nil
}
