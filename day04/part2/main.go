package main

import (
	"flag"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")

	pidRE = regexp.MustCompile(`^\d{9}$`)
	rgbRE = regexp.MustCompile(`^#[0-9a-f]{6}$`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	passports := readPassports(filename)
	log.Printf("%v valid passports", len(passports))
}

type Passport struct {
	valid    bool
	required int
	pairs    map[string]string
}

var required = map[string]bool{
	"byr": true, // (Birth Year)
	"iyr": true, // (Issue Year)
	"eyr": true, // (Expiration Year)
	"hgt": true, // (Height)
	"hcl": true, // (Hair Color)
	"ecl": true, // (Eye Color)
	"pid": true, // (Passport ID)
	// "cid" // (Country ID)
}

func readPassports(filename string) []*Passport {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	var result []*Passport
	passports := strings.Split(string(buf), "\n\n")
	for _, passport := range passports {
		passport = strings.ReplaceAll(passport, "\n", " ")
		passport = strings.TrimSpace(passport)
		if passport == "" {
			continue
		}

		parts := strings.Split(passport, " ")
		p := &Passport{pairs: map[string]string{}}
		for _, part := range parts {
			keyVal := strings.Split(part, ":")
			if len(keyVal) != 2 {
				log.Fatalf("unexpected passport: %v", passport)
			}
			p.pairs[keyVal[0]] = keyVal[1]
			if validPair(keyVal[0], keyVal[1]) {
				p.required++
			}
		}

		if p.required == len(required) {
			p.valid = true
			result = append(result, p)
		}
	}

	return result
}

func validPair(key, val string) bool {
	if required[key] {
		switch key {
		case "byr":
			// byr (Birth Year) - four digits; at least 1920 and at most 2002.
			if validNum(val, 1920, 2002) {
				return true
			}
		case "iyr":
			// iyr (Issue Year) - four digits; at least 2010 and at most 2020.
			if validNum(val, 2010, 2020) {
				return true
			}
		case "eyr":
			// eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
			if validNum(val, 2020, 2030) {
				return true
			}
		case "hgt":
			// hgt (Height) - a number followed by either cm or in:
			if strings.HasSuffix(val, "cm") && validNum(val[:len(val)-2], 150, 193) {
				//   If cm, the number must be at least 150 and at most 193.
				return true
			} else if strings.HasSuffix(val, "in") && validNum(val[:len(val)-2], 59, 76) {
				//   If in, the number must be at least 59 and at most 76.
				return true
			}
		case "hcl":
			// hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
			if validRGB(val) {
				return true
			}
		case "ecl":
			// ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
			switch val {
			case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
				return true
			}
		case "pid":
			// pid (Passport ID) - a nine-digit number, including leading zeroes.
			if validPID(val) {
				return true
			}
			// cid (Country ID) - ignored, missing or not.
		}
	}
	return false
}

func validNum(s string, min, max int) bool {
	v, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return v >= min && v <= max
}

func validRGB(s string) bool {
	return rgbRE.MatchString(s)
}

func validPID(s string) bool {
	return pidRE.MatchString(s)
}

func check(fmtStr string, args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		log.Fatalf(fmtStr, args...)
	}
}

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}
