package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const inputPrefix = "input"

func main() {
	log.SetFlags(0)
	foods, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	ingsMaybeAllergens := make(map[string]struct{})
	allergenCandidates := make(map[string]map[string]struct{})
	for _, all := range allAlergens(foods) {
		ings := ingredientsMaybeContainingAllergen(foods, all)
		allergenCandidates[all] = make(map[string]struct{})
		for _, ing := range ings {
			allergenCandidates[all][ing] = struct{}{}
		}
		for _, ing := range ings {
			ingsMaybeAllergens[ing] = struct{}{}
		}
	}

	var ingsWithoutAllergens []string
	for _, ing := range allIngredients(foods) {
		if _, ok := ingsMaybeAllergens[ing]; !ok {
			ingsWithoutAllergens = append(ingsWithoutAllergens, ing)
		}
	}
	part1 := countIngredients(foods, ingsWithoutAllergens)

	dangerousIngs := determineDangerousIngredients(foods, allergenCandidates)
	part2 := strings.Join(dangerousIngs, ",")

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %s\n", part2)
}

func determineDangerousIngredients(foods []*food, allergenCandidates map[string]map[string]struct{}) []string {
	allergenIngredient := make(map[string]string)
	var queue []string
	for all := range allergenCandidates {
		queue = append(queue, all)
	}

	for len(queue) > 0 {
		var all string
		all, queue = queue[0], queue[1:]
		ings := allergenCandidates[all]
		if len(ings) == 0 {
			panic("bad state")
		}
		if len(ings) > 1 {
			queue = append(queue, all)
			continue
		}
		var ing string
		for k := range ings {
			ing = k
		}
		allergenIngredient[all] = ing
		for _, ings := range allergenCandidates {
			delete(ings, ing)
		}
	}

	var alls []string
	for all := range allergenIngredient {
		alls = append(alls, all)
	}
	sort.Strings(alls)
	ings := make([]string, len(alls))
	for i, all := range alls {
		ings[i] = allergenIngredient[all]
	}
	return ings
}

func countIngredients(foods []*food, ings []string) int {
	n := 0
	for _, f := range foods {
	ingredients:
		for _, ing := range f.ingredients {
			for _, i := range ings {
				if ing == i {
					n++
					continue ingredients
				}
			}
		}
	}
	return n
}

func ingredientsMaybeContainingAllergen(foods []*food, allergen string) []string {
	var foodsContaining []*food
foods:
	for _, f := range foods {
		for _, all := range f.allergens {
			if all == allergen {
				foodsContaining = append(foodsContaining, f)
				continue foods
			}
		}
	}
	ingredientCounts := make(map[string]int)
	for _, f := range foodsContaining {
		for _, ing := range f.ingredients {
			ingredientCounts[ing]++
		}
	}
	var ings []string
	for ing, count := range ingredientCounts {
		if count == len(foodsContaining) {
			ings = append(ings, ing)
		}
	}
	return ings
}

func allAlergens(foods []*food) []string {
	m := make(map[string]struct{})
	for _, f := range foods {
		for _, all := range f.allergens {
			m[all] = struct{}{}
		}
	}
	var alls []string
	for a := range m {
		alls = append(alls, a)
	}
	sort.Strings(alls)
	return alls
}

func allIngredients(foods []*food) []string {
	m := make(map[string]struct{})
	for _, f := range foods {
		for _, ing := range f.ingredients {
			m[ing] = struct{}{}
		}
	}
	var ings []string
	for ing := range m {
		ings = append(ings, ing)
	}
	sort.Strings(ings)
	return ings
}

func readInput(filename string) ([]*food, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var foods []*food
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.Index(line, "(")
		j := strings.Index(line, ")")
		f := &food{
			ingredients: strings.Fields(line[:i]),
		}
		for _, all := range strings.Fields(line[i+1 : j])[1:] {
			f.allergens = append(f.allergens, strings.Trim(all, ","))
		}
		foods = append(foods, f)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return foods, nil
}

type food struct {
	ingredients []string
	allergens   []string
}
