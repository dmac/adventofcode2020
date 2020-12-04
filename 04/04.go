package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	passports, err := readInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	p1Valid := 0
	p2Valid := 0
	for _, pp := range passports {
		if pp.byr == "" || pp.iyr == "" || pp.eyr == "" || pp.hgt == "" ||
			pp.hcl == "" || pp.ecl == "" || pp.pid == "" {
			continue
		}
		p1Valid++

		byr, err := strconv.Atoi(pp.byr)
		if err != nil {
			continue
		}
		if byr < 1920 || byr > 2002 {
			continue
		}

		iyr, err := strconv.Atoi(pp.iyr)
		if err != nil {
			continue
		}
		if iyr < 2010 || iyr > 2020 {
			continue
		}

		eyr, err := strconv.Atoi(pp.eyr)
		if err != nil {
			continue
		}
		if eyr < 2020 || eyr > 2030 {
			continue
		}

		if !strings.HasSuffix(pp.hgt, "cm") && !strings.HasSuffix(pp.hgt, "in") {
			continue
		}
		hgt, err := strconv.Atoi(pp.hgt[:len(pp.hgt)-2])
		if err != nil {
			continue
		}
		if strings.HasSuffix(pp.hgt, "cm") {
			if hgt < 150 || hgt > 193 {
				continue
			}
		} else if strings.HasSuffix(pp.hgt, "in") {
			if hgt < 59 || hgt > 76 {
				continue
			}
		}

		if pp.hcl[0] != '#' {
			continue
		}
		for _, c := range pp.hcl[1:] {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				continue
			}
		}

		if pp.ecl != "amb" && pp.ecl != "blu" && pp.ecl != "brn" && pp.ecl != "gry" &&
			pp.ecl != "grn" && pp.ecl != "hzl" && pp.ecl != "oth" {
			continue
		}

		if len(pp.pid) != 9 {
			continue
		}
		if _, err := strconv.Atoi(pp.pid); err != nil {
			continue
		}

		p2Valid++
	}
	fmt.Printf("Part 1: %d\n", p1Valid)
	fmt.Printf("Part 2: %d\n", p2Valid)
}

type passport struct {
	byr string
	iyr string
	eyr string
	hgt string
	hcl string
	ecl string
	pid string
	cid string
}

func readInput(filename string) ([]*passport, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var passports []*passport
	pp := new(passport)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			passports = append(passports, pp)
			pp = new(passport)
			continue
		}
		for _, field := range strings.Fields(line) {
			pair := strings.Split(field, ":")
			switch pair[0] {
			case "byr":
				pp.byr = pair[1]
			case "iyr":
				pp.iyr = pair[1]
			case "eyr":
				pp.eyr = pair[1]
			case "hgt":
				pp.hgt = pair[1]
			case "hcl":
				pp.hcl = pair[1]
			case "ecl":
				pp.ecl = pair[1]
			case "pid":
				pp.pid = pair[1]
			case "cid":
				pp.cid = pair[1]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return passports, nil
}
