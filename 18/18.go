package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	exprs, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	part1 := 0
	for _, expr := range exprs {
		part1 += expr.evalOrder()
	}

	part2 := 0
	for _, expr := range exprs {
		part2 += expr.evalPrecedence()
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

const (
	typeNumber = iota
	typeOperator
	typeExpression
)

type expression struct {
	typ int

	num int
	op  byte

	children []*expression
}

func (e *expression) evalOrder() int {
	switch e.typ {
	case typeNumber:
		return e.num
	case typeOperator:
		panic("cannot evaluate operator")
	case typeExpression:
		var ans int
		var op byte
		for _, expr := range e.children {
			var val int
			switch expr.typ {
			case typeNumber, typeExpression:
				val = expr.evalOrder()
			case typeOperator:
				op = expr.op
				continue
			}
			switch op {
			case 0:
				ans = val
			case '+':
				ans += val
			case '*':
				ans *= val
			default:
				panic(fmt.Sprintf("unknown op %q", e.children[1].op))
			}
		}
		return ans
	default:
		panic(fmt.Sprintf("unknown expression type %d", e.typ))
	}
}

func (e *expression) evalPrecedence() int {
	switch e.typ {
	case typeNumber:
		return e.num
	case typeOperator:
		panic("cannot evaluate operator")
	case typeExpression:
		var firstPass []*expression
		var op byte
		var ans int
		for _, expr := range e.children {
			var val int
			switch expr.typ {
			case typeNumber, typeExpression:
				val = expr.evalPrecedence()
			case typeOperator:
				op = expr.op
				continue
			}
			switch op {
			case 0:
				ans = val
			case '+':
				ans += val
			case '*':
				expr := &expression{typ: typeNumber, num: ans}
				firstPass = append(firstPass, expr)
				expr = &expression{typ: typeOperator, op: '*'}
				firstPass = append(firstPass, expr)
				ans = val
				op = 0
			default:
				panic(fmt.Sprintf("unknown op %q", e.children[1].op))
			}
		}
		expr := &expression{typ: typeNumber, num: ans}
		firstPass = append(firstPass, expr)

		op = 0
		ans = 0
		for _, expr := range firstPass {
			var val int
			switch expr.typ {
			case typeNumber, typeExpression:
				val = expr.evalPrecedence()
			case typeOperator:
				op = expr.op
				continue
			}
			switch op {
			case 0:
				ans = val
			case '+':
				panic("unexpected + on second pass")
			case '*':
				ans *= val
			default:
				panic(fmt.Sprintf("unknown op %q", e.children[1].op))
			}
		}

		return ans
	default:
		panic(fmt.Sprintf("unknown expression type %d", e.typ))
	}
}

func parseExpression(s string) (*expression, error) {
	expr := &expression{typ: typeExpression}
	for {
		var ex *expression
		var err error
		ex, s, err = parseNextExpression(s)
		if err != nil {
			return nil, err
		}
		if ex == nil && s == "" {
			break
		}
		expr.children = append(expr.children, ex)
	}
	if len(expr.children) == 1 {
		return expr.children[0], nil
	}
	return expr, nil
}

func parseNextExpression(s string) (*expression, string, error) {
	s = discardWhitespace(s)
	if s == "" {
		return nil, "", nil
	}
	if unicode.IsDigit(rune(s[0])) {
		i := scanPastNextSpace(s)
		n, err := strconv.Atoi(strings.TrimSpace(s[:i]))
		if err != nil {
			return nil, "", err
		}
		expr := &expression{
			typ: typeNumber,
			num: n,
		}
		return expr, s[i:], nil
	}
	if s[0] == '+' || s[0] == '*' {
		expr := &expression{
			typ: typeOperator,
			op:  s[0],
		}
		return expr, s[1:], nil
	}
	if s[0] == '(' {
		i, err := scanToBalancedParen(s)
		if err != nil {
			return nil, "", err
		}
		expr, err := parseExpression(s[1:i])
		if err != nil {
			return nil, "", err
		}
		return expr, s[i+1:], nil
	}
	return nil, s, fmt.Errorf("unable to parse %q", s)
}

func discardWhitespace(s string) string {
	for i := 0; i < len(s); i++ {
		if !unicode.IsSpace(rune(s[i])) {
			return s[i:]
		}
	}
	return s
}

func scanPastNextSpace(s string) int {
	space := false
	i := 0
	for ; i < len(s); i++ {
		if unicode.IsSpace(rune(s[i])) {
			space = true
		} else if space {
			return i
		}
	}
	return i
}

func scanToBalancedParen(s string) (int, error) {
	parens := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			parens++
		} else if s[i] == ')' {
			parens--
		}
		if parens == 0 {
			return i, nil
		}
	}
	return 0, errors.New("unbalanced parens")
}

func readInput(filename string) ([]*expression, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var exprs []*expression
	for scanner.Scan() {
		expr, err := parseExpression(scanner.Text())
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return exprs, nil
}
