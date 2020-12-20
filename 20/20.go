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
	tiles, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	size := int(math.Sqrt(float64(len(tiles))))
	table := make([][]*tile, size)
	for i := range table {
		table[i] = make([]*tile, size)
	}
	if !assemble(table, tiles) {
		log.Fatal("no solution found")
	}

	part1 := table[0][0].id
	part1 *= table[0][len(table[0])-1].id
	part1 *= table[len(table)-1][0].id
	part1 *= table[len(table)-1][len(table[0])-1].id

	part2 := 0
	image := assembleImage(table)
	for _, img0 := range []*tile{
		image,
		image.flipH(),
		image.flipV(),
	} {
		for rot := 0; rot < 4; rot++ {
			img1 := img0.rotateCW(rot)
			if countMonsters(img1) > 0 {
				part2 = countNonMonsters(img1)
				goto done
			}
		}
	}
done:
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func countNonMonsters(t *tile) int {
	count := 0
	for i := 0; i < len(t.pixels); i++ {
		for j := 0; j < len(t.pixels[0]); j++ {
			if t.pixels[i][j] == '#' {
				count++
			}
		}
	}
	return count
}

func countMonsters(t *tile) int {
	mask := [][]bool{
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false},
		{true, false, false, false, false, true, true, false, false, false, false, true, true, false, false, false, false, true, true, true},
		{false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, false},
	}
	count := 0
	for r := 0; r < len(t.pixels)-len(mask)+1; r++ {
		for c := 0; c < len(t.pixels[0])-len(mask[0])+1; c++ {
			found := true
		loop:
			for mr := 0; mr < len(mask); mr++ {
				for mc := 0; mc < len(mask[0]); mc++ {
					if mask[mr][mc] && t.pixels[r+mr][c+mc] != '#' {
						found = false
						break loop
					}
				}
			}
			if found {
				count++
				for mr := 0; mr < len(mask); mr++ {
					for mc := 0; mc < len(mask[0]); mc++ {
						if mask[mr][mc] {
							t.pixels[r+mr][c+mc] = 'O'
						}
					}
				}
			}
		}
	}
	return count
}

func assembleImage(table [][]*tile) *tile {
	for i := range table {
		for _, t := range table[i] {
			t.pixels = t.pixels[1 : len(t.pixels)-1]
			for j, row := range t.pixels {
				t.pixels[j] = row[1 : len(row)-1]
			}
		}
	}
	size := (len(table[0][0].pixels)) * len(table)
	t := &tile{pixels: make([][]byte, size)}
	for i := range t.pixels {
		t.pixels[i] = make([]byte, size)
		tr := i / len(table[0][0].pixels)
		pr := i % len(table[0][0].pixels)
		for tc := 0; tc < len(table[0]); tc++ {
			pixels := table[tr][tc].pixels
			start := tc * len(pixels)
			end := start + len(pixels)
			copy(t.pixels[i][start:end], table[tr][tc].pixels[pr])
		}
	}
	return t
}

func assemble(table [][]*tile, tiles []*tile) bool {
	size := 2*len(table) - 1
	bigTable := make([][]*tile, size)
	for i := range bigTable {
		bigTable[i] = make([]*tile, size)
	}
	mid := size / 2
	bigTable[mid][mid] = tiles[0]
	if !assembleHelper(bigTable, tiles[1:]) {
		return false
	}
	var r, c int
	for r = 0; r < len(bigTable); r++ {
		for c = 0; c < len(bigTable[0]); c++ {
			if bigTable[r][c] != nil {
				goto done
			}
		}
	}
done:
	for row := 0; row < len(table); row++ {
		for col := 0; col < len(table[0]); col++ {
			table[row][col] = bigTable[row+r][col+c]
		}
	}
	return true
}

func assembleHelper(table [][]*tile, tiles []*tile) bool {
	if len(tiles) == 0 {
		return true
	}
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[0]); c++ {
			if table[r][c] == nil {
				continue
			}
			for i := range tiles {
				newTiles := append(append([]*tile{}, tiles[:i]...), tiles[i+1:]...)
				for _, t := range []*tile{
					tiles[i],
					tiles[i].flipH(),
					tiles[i].flipV(),
				} {
					for rot := 0; rot < 4; rot++ {
						if r > 0 && table[r-1][c] == nil {
							table[r-1][c] = t.rotateCW(rot)
							if checkValid(table, r-1, c) {
								if assembleHelper(table, newTiles) {
									return true
								}
							}
							table[r-1][c] = nil
						}
						if r < len(table)-1 && table[r+1][c] == nil {
							table[r+1][c] = t.rotateCW(rot)
							if checkValid(table, r+1, c) {
								if assembleHelper(table, newTiles) {
									return true
								}
							}
							table[r+1][c] = nil
						}
						if c > 0 && table[r][c-1] == nil {
							table[r][c-1] = t.rotateCW(rot)
							if checkValid(table, r, c-1) {
								if assembleHelper(table, newTiles) {
									return true
								}
							}
							table[r][c-1] = nil
						}
						if c < len(table[0])-1 && table[r][c+1] == nil {
							table[r][c+1] = t.rotateCW(rot)
							if checkValid(table, r, c+1) {
								if assembleHelper(table, newTiles) {
									return true
								}
							}
							table[r][c+1] = nil
						}
					}
				}
			}
		}
	}
	return false
}

func checkValid(table [][]*tile, row, col int) bool {
	t := table[row][col]
	if row > 0 {
		top := table[row-1][col]
		if top != nil && top.edge(Bottom) != t.edge(Top) {
			return false
		}
	}
	if row < len(table)-1 {
		bot := table[row+1][col]
		if bot != nil && bot.edge(Top) != t.edge(Bottom) {
			return false
		}
	}
	if col > 0 {
		left := table[row][col-1]
		if left != nil && left.edge(Right) != t.edge(Left) {
			return false
		}
	}
	if col < len(table[0])-1 {
		right := table[row][col+1]
		if right != nil && right.edge(Left) != t.edge(Right) {
			return false
		}
	}
	return true
}

func readInput(filename string) ([]*tile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var tiles []*tile
	var t *tile
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			tiles = append(tiles, t)
			continue
		}
		if strings.Contains(line, "Tile") {
			i := strings.Index(line, " ")
			j := strings.Index(line, ":")
			id, err := strconv.Atoi(line[i+1 : j])
			if err != nil {
				return nil, err
			}
			t = &tile{id: id}
			continue
		}
		t.pixels = append(t.pixels, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	tiles = append(tiles, t)
	return tiles, nil
}

type tile struct {
	id     int
	pixels [][]byte

	fh  bool
	fv  bool
	rot int
}

func (t *tile) String() string {
	var sb strings.Builder
	for _, line := range t.pixels {
		sb.WriteString(string(line))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *tile) rotateCW(times int) *tile {
	if times == 0 {
		return t
	}
	if len(t.pixels) != len(t.pixels[0]) {
		fmt.Println(t)
		panic(fmt.Sprintf("expected square tile, got %dx%d", len(t.pixels), len(t.pixels[0])))
	}
	size := len(t.pixels)
	tt := &tile{
		id:     t.id,
		pixels: make([][]byte, size),
		fh:     t.fh,
		fv:     t.fv,
		rot:    (t.rot + 1) % 4,
	}
	for i, line := range t.pixels {
		tt.pixels[i] = make([]byte, len(line))
		copy(tt.pixels[i], line)
	}
	for r := 0; r < len(t.pixels); r++ {
		for c := 0; c < len(t.pixels); c++ {
			tt.pixels[r][c] = t.pixels[size-c-1][r]
		}
	}
	return tt.rotateCW(times - 1)
}

func (t *tile) flipH() *tile {
	if len(t.pixels) != len(t.pixels[0]) {
		fmt.Println(t)
		panic(fmt.Sprintf("expected square tile, got %dx%d", len(t.pixels), len(t.pixels[0])))
	}
	size := len(t.pixels)
	tt := &tile{
		id:     t.id,
		pixels: make([][]byte, size),
		fh:     !t.fh,
		fv:     t.fv,
		rot:    t.rot,
	}
	for i, line := range t.pixels {
		tt.pixels[i] = make([]byte, len(line))
		for j := 0; j < len(tt.pixels[i]); j++ {
			tt.pixels[i][j] = line[len(line)-1-j]
		}
	}
	return tt
}

func (t *tile) flipV() *tile {
	if len(t.pixels) != len(t.pixels[0]) {
		fmt.Println(t)
		panic(fmt.Sprintf("expected square tile, got %dx%d", len(t.pixels), len(t.pixels[0])))
	}
	size := len(t.pixels)
	tt := &tile{
		id:     t.id,
		pixels: make([][]byte, size),
		fh:     t.fh,
		fv:     !t.fv,
		rot:    t.rot,
	}
	for i, line := range t.pixels {
		tt.pixels[i] = make([]byte, len(line))
	}
	for r := 0; r < len(t.pixels); r++ {
		for c := 0; c < len(t.pixels[0]); c++ {
			tt.pixels[r][c] = t.pixels[len(t.pixels)-r-1][c]
		}
	}
	return tt
}

const (
	Top = iota
	Right
	Bottom
	Left
)

func (t *tile) edge(side int) string {
	switch side {
	case Top:
		return string(t.pixels[0])
	case Bottom:
		return string(t.pixels[len(t.pixels)-1])
	case Left:
		b := make([]byte, len(t.pixels))
		for r := 0; r < len(t.pixels); r++ {
			b[r] = t.pixels[r][0]
		}
		return string(b)
	case Right:
		b := make([]byte, len(t.pixels))
		for r := 0; r < len(t.pixels); r++ {
			b[r] = t.pixels[r][len(t.pixels[r])-1]
		}
		return string(b)
	default:
		panic(fmt.Sprintf("unexpected side: %d", side))
	}
}

func stateHash(table [][]*tile, tiles []*tile) string {
	var sb strings.Builder
	for r, row := range table {
		for c, t := range row {
			if t == nil {
				sb.WriteString(fmt.Sprintf("%d,%d:nil\n", r, c))
				continue
			}
			sb.WriteString(fmt.Sprintf("%d,%d:%d_%d", r, c, t.id, t.rot))
			if t.fh {
				sb.WriteString("H")
			}
			if t.fv {
				sb.WriteString("V")
			}
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
	for i, t := range tiles {
		sb.WriteString(fmt.Sprintf("%d:%d\n", i, t.id))
	}
	return sb.String()
}

