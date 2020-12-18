// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"regexp"
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
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	lines := strings.Split(string(buf), "\n")
	puz := &Puzzle{space: map[key]bool{}}
	var y int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for x := 0; x < len(line); x++ {
			if line[x:x+1] == "#" {
				puz.space[key{x: x, y: y}] = true
			}
		}
		y++
	}

	for cycle := 1; cycle <= 6; cycle++ {
		after := puz.cycle()
		puz.space = after
	}

	log.Printf("Solution: %v", len(puz.space))
}

type key struct {
	x, y, z int
}

func (p *Puzzle) cycle() map[key]bool {
	r := map[key]bool{}
	inactiveNeighbors := map[key]bool{}
	for k := range p.space {
		activeNeighbors := map[key]bool{}
		p.neighbors(k, true, activeNeighbors)
		if len(activeNeighbors) == 2 || len(activeNeighbors) == 3 {
			r[k] = true
		}
		p.neighbors(k, false, inactiveNeighbors)
	}

	for k := range inactiveNeighbors {
		activeNeighbors := map[key]bool{}
		p.neighbors(k, true, activeNeighbors)
		if len(activeNeighbors) == 3 {
			r[k] = true
		}
	}

	return r
}

func (p *Puzzle) neighbors(k key, active bool, result map[key]bool) {
	for m := -1; m <= 1; m++ {
		for j := -1; j <= 1; j++ {
			for i := -1; i <= 1; i++ {
				if i == 0 && j == 0 && m == 0 {
					continue
				}
				nk := key{x: k.x + i, y: k.y + j, z: k.z + m}
				if p.space[nk] == active {
					result[nk] = true
				}
			}
		}
	}
}

type Puzzle struct {
	space map[key]bool
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
