package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const inputPrefix = "input"

func main() {
	rules, messages, err := readInput(inputPrefix + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	expanded := expandRules(rules)

	valid := make(map[string]struct{})
	for _, ex := range expanded[0].expanded {
		valid[ex] = struct{}{}
	}

	part1 := 0
	for _, msg := range messages {
		if _, ok := valid[msg]; ok {
			part1++
		}
	}

	rules, messages, err = readInput(inputPrefix + "2.txt")
	if err != nil {
		log.Fatal(err)
	}
	expanded = expandRules(rules)
	// printAllExpandedRules(expanded)
	// fmt.Println(expanded[31])
	// fmt.Println(expanded[42])
	part2 := 0
	for _, msg := range messages {
		if isValidForPart2(expanded, msg) {
			part2++
		}
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

// To be valid for a part two, a message must start with any number of
// repeating 42 strings (at least 1), then have some number of 42 strings
// followed by the same number of 31 strings (at least 1 each).
func isValidForPart2(expanded map[int]*rule, msg string) bool {
	thirtyone := expanded[31]
	fortytwo := expanded[42]
	found := false
restart:
	if found && satisfiesSecondHalf(thirtyone, fortytwo, msg) {
		return true
	}
	for _, pre := range fortytwo.expanded {
		if strings.HasPrefix(msg, pre) {
			found = true
			msg = msg[len(pre):]
			goto restart
		}
	}
	return false
}

func satisfiesSecondHalf(thirtyone, fortytwo *rule, msg string) bool {
	found := false
restart:
	for _, pre := range fortytwo.expanded {
		if strings.HasPrefix(msg, pre) {
			for _, suf := range thirtyone.expanded {
				if strings.HasSuffix(msg, suf) {
					found = true
					msg = msg[len(pre) : len(msg)-len(suf)]
					goto restart
				}
			}
		}
	}
	return found && len(msg) == 0
}

func printAllExpandedRules(expanded map[int]*rule) {
	var rs []*rule
	for _, r := range expanded {
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].id < rs[j].id
	})
	for _, r := range rs {
		fmt.Println(r)
	}
}

func expandRules(rules map[int]*rule) map[int]*rule {
	expanded := make(map[int]*rule)
	var queue []*rule
	for _, r := range rules {
		if r.literal != "" {
			r.expanded = []string{r.literal}
			expanded[r.id] = r
		} else {
			queue = append(queue, r)
		}
	}
	seen := make(map[string]struct{})
	var curr *rule
	for len(queue) > 0 {
		h := hashQueue(queue)
		if _, ok := seen[h]; ok {
			break
		}
		seen[h] = struct{}{}
		curr, queue = queue[0], queue[1:]
		if curr.literal != "" {
			panic("literal rule found in expansion queue")
		}
		allExpanded := true
		for _, ids := range curr.rules {
			for _, id := range ids {
				if _, ok := expanded[id]; !ok {
					allExpanded = false
				}
			}
		}
		if !allExpanded {
			queue = append(queue, curr)
			continue
		}
		for _, ids := range curr.rules {
			switch len(ids) {
			case 1:
				r0 := expanded[ids[0]]
				for _, e0 := range r0.expanded {
					curr.expanded = append(curr.expanded, e0)
				}
			case 2:
				r0 := expanded[ids[0]]
				r1 := expanded[ids[1]]
				for _, e0 := range r0.expanded {
					for _, e1 := range r1.expanded {
						curr.expanded = append(curr.expanded, e0+e1)
					}
				}
			case 3:
				r0 := expanded[ids[0]]
				r1 := expanded[ids[1]]
				r2 := expanded[ids[2]]
				for _, e0 := range r0.expanded {
					for _, e1 := range r1.expanded {
						for _, e2 := range r2.expanded {
							curr.expanded = append(curr.expanded, e0+e1+e2)
						}
					}
				}
			default:
				panic(fmt.Sprintf("doesn't handle len(ids) == %d", len(ids)))
			}
		}
		expanded[curr.id] = curr
	}
	return expanded
}

func hashQueue(queue []*rule) string {
	return fmt.Sprintf("%#v\n", queue)
}

func parseRule(s string) (*rule, error) {
	parts := strings.Split(s, ": ")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	r := &rule{id: id}
	rest := strings.TrimSpace(parts[1])
	if rest[0] == '"' {
		r.literal = rest[1 : len(rest)-1]
		return r, nil
	}
	for _, conjunction := range strings.Split(rest, "|") {
		var ids []int
		for _, sid := range strings.Fields(conjunction) {
			id, err := strconv.Atoi(sid)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		r.rules = append(r.rules, ids)
	}
	return r, nil
}

func readInput(filename string) (map[int]*rule, []string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	rules := make(map[int]*rule)
	var messages []string
	doneWithRules := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		if line == "" {
			doneWithRules = true
			continue
		}
		if doneWithRules {
			messages = append(messages, line)
			continue
		}
		rule, err := parseRule(line)
		if err != nil {
			return nil, nil, err
		}
		rules[rule.id] = rule
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return rules, messages, nil
}

type rule struct {
	id       int
	literal  string
	rules    [][]int
	expanded []string
}

func (r *rule) String() string {
	if r.literal != "" {
		return fmt.Sprintf("%d:%s", r.id, r.literal)
	}
	if r.expanded != nil {
		return fmt.Sprintf("%d:%s", r.id, r.expanded)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d:", r.id))
	for i, ids := range r.rules {
		if i > 0 {
			sb.WriteString("|")
		}
		sb.WriteString("[")
		for j, id := range ids {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("%d", id))
		}
		sb.WriteString("]")
	}
	return sb.String()
}
