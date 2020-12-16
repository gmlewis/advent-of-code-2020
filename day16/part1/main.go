// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

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

	valsRE = regexp.MustCompile(`^(.*): (\d+)-(\d+) or (\d+)-(\d+)$`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	groups := strings.Split(string(buf), "\n\n")
	puz := &Puzzle{fields: map[string]valRanges{}}
	puz.processRanges(groups[0])
	errorRate := puz.errorRate(groups[2])

	log.Printf("Solution: %v", errorRate)
}

type Puzzle struct {
	fields map[string]valRanges
}

type valRanges struct {
	s1 int
	e1 int
	s2 int
	e2 int
}

func (p *Puzzle) processRanges(buf string) {
	for _, line := range strings.Split(buf, "\n") {
		m := valsRE.FindStringSubmatch(line)
		if len(m) != 6 {
			log.Fatalf("bad input/regexp: %v", line)
		}
		s1, err := strconv.Atoi(m[2])
		check("s1: %v", err)
		e1, err := strconv.Atoi(m[3])
		check("e1: %v", err)
		s2, err := strconv.Atoi(m[4])
		check("s2: %v", err)
		e2, err := strconv.Atoi(m[5])
		check("e2: %v", err)
		p.fields[m[1]] = valRanges{s1: s1, e1: e1, s2: s2, e2: e2}
	}
}

func (p *Puzzle) errorRate(buf string) int {
	var result int
	for _, line := range strings.Split(buf, "\n")[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		for _, part := range parts {
			v, err := strconv.Atoi(part)
			check("v: %v - %v", line, err)
			if !p.isValid(v) {
				result += v
			}
		}
	}
	return result
}

func (p *Puzzle) isValid(n int) bool {
	for _, v := range p.fields {
		if (n >= v.s1 && n <= v.e1) || (n >= v.s2 && n <= v.e2) {
			return true
		}
	}
	return false
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
