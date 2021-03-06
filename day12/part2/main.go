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
			p.wnpos += amt
		case "S":
			p.wnpos -= amt
		case "E":
			p.wepos += amt
		case "W":
			p.wepos -= amt
		case "L":
			for amt > 0 {
				p.wepos, p.wnpos = -p.wnpos, p.wepos
				amt -= 90
			}
		case "R":
			for amt > 0 {
				p.wepos, p.wnpos = p.wnpos, -p.wepos
				amt -= 90
			}
		case "F":
			p.npos += amt * p.wnpos
			p.epos += amt * p.wepos
		}
		// log.Printf("%v: (%v,%v) (%v,%v)", line, p.epos, p.npos, p.wepos, p.wnpos)
	}
}

type Puzzle struct {
	wepos int
	wnpos int
	epos  int
	npos  int
	lines []string
}

func readPuzzle(filename string) *Puzzle {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	lines := strings.Split(string(buf), "\n")
	puz := &Puzzle{lines: lines, wepos: 10, wnpos: 1}

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
