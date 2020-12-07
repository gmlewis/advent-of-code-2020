// -*- compile-command: "go run main.go ../example1.txt ../example2.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")

	lineRE     = regexp.MustCompile(`^(.*?) bags contain (.*)\.$`)
	containsRE = regexp.MustCompile(`,?\s*(\d+) (.*?) bags?`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

type Contains struct {
	quant int
	color string
}

func process(filename string) {
	rules := map[string][]*Contains{}

	lines := readLines(filename)
	// log.Printf("lines=%v", lines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		m := lineRE.FindStringSubmatch(line)
		if len(m) != 3 {
			log.Fatalf("bad parse: %v", line)
		}
		// log.Printf("len(m)=%v: %v\n%#v", len(m), line, m)
		if m[2] == "no other bags" {
			rules[m[1]] = nil
			continue
		}

		m2 := containsRE.FindAllStringSubmatch(m[2], -1)
		// log.Printf("len(m2)=%v: %v\n%#v", len(m2), m[2], m2)
		for _, f := range m2 {
			q, err := strconv.Atoi(f[1])
			if err != nil {
				log.Fatalf("strconv.Atoi(%q): %v", f[1], err)
			}
			rules[m[1]] = append(rules[m[1]], &Contains{quant: q, color: f[2]})
		}
	}

	count := contains("shiny gold", rules, nil)

	fmt.Printf("Solution: %v\n", count-1)
}

func contains(key string, rules map[string][]*Contains, seen map[string]bool) int {
	total := 1
	for _, v := range rules[key] {
		count := contains(v.color, rules, seen)
		total += (count * v.quant)
	}
	// log.Printf("contains(%q)=%v", key, total)
	return total
}

func readLines(filename string) []string {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)
	lines := strings.Split(string(buf), "\n")
	return lines
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
