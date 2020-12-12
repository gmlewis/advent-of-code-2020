// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
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
	puz := readPuzzle(filename)

	puz.iterate()

	dist := puz.manhattan()

	log.Printf("Solution: %v", dist)
}

func (p *Puzzle) manhattan() int {
	x := p.epos
	if x < 0 {
		x = -x
	}
	y := p.npos
	if y < 0 {
		y = -y
	}
	return x + y
}

func (p *Puzzle) iterate() {
	for _, line := range p.lines {
		if line == "" {
			continue
		}

		amt, err := strconv.Atoi(line[1:])
		check("Atoi: %v", err)

		switch line[0:1] {
		case "N":
			p.npos += amt
		case "S":
			p.npos -= amt
		case "E":
			p.epos += amt
		case "W":
			p.epos -= amt
		case "L":
			for amt > 0 {
				p.de, p.dn = -p.dn, p.de
				amt -= 90
			}
		case "R":
			for amt > 0 {
				p.de, p.dn = p.dn, -p.de
				amt -= 90
			}
		case "F":
			p.npos += amt * p.dn
			p.epos += amt * p.de
		}
		// log.Printf("%v: (%v,%v) (%v,%v)", line, p.epos, p.npos, p.de, p.dn)
	}
}

type Puzzle struct {
	de    int
	dn    int
	epos  int
	npos  int
	lines []string
}

func readPuzzle(filename string) *Puzzle {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	lines := strings.Split(string(buf), "\n")
	puz := &Puzzle{lines: lines, de: 1}

	return puz
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
