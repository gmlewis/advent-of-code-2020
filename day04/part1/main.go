package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")
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
			if required[keyVal[0]] {
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
