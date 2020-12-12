// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
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
	puz := readPuzzle(filename)

	for puz.iterate() {
	}

	occupied := puz.occupied()

	log.Printf("Solution: %v", occupied)
}

func (p *Puzzle) occupied() int {
	var count int
	for _, v := range p.grid {
		if v == "#" {
			count++
		}
	}
	return count
}

func (p *Puzzle) iterate() bool {
	r := map[string]string{}
	var changesMade bool
	for k, v := range p.grid {
		adj := p.countAdjacentOccupied(k)
		// log.Printf("k=%v, v=%v, adj=%v", k, v, adj)
		if v == "L" && adj == 0 {
			r[k] = "#"
			changesMade = true
		} else if v == "#" && adj >= 5 {
			r[k] = "L"
			changesMade = true
		} else {
			r[k] = v
		}
	}

	if changesMade {
		p.grid = r
	}

	return changesMade
}

func (p *Puzzle) countAdjacentOccupied(k string) int {
	x, y := p.parseKey(k)
	var adj int

	t := func(dx, dy int) {
		u, v := x+dx, y+dy
		for u >= 0 && v >= 0 && u < p.width && v < p.height {
			val := p.lookup(u, v)
			if val == "#" {
				adj++
				return
			}
			if val == "L" {
				return
			}
			u += dx
			v += dy
		}
	}

	t(-1, -1)
	t(0, -1)
	t(+1, -1)
	t(-1, 0)

	t(+1, 0)
	t(-1, +1)
	t(0, +1)
	t(+1, +1)

	return adj
}

type Puzzle struct {
	width  int
	height int
	grid   map[string]string
	keyX   map[string]int
	keyY   map[string]int
}

func (p *Puzzle) String() string {
	var s string
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			v := p.lookup(x, y)
			if v == "" {
				s += "."
			} else {
				s += v
			}
		}
		s += "\n"
	}
	return s
}

func readPuzzle(filename string) *Puzzle {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	puz := &Puzzle{grid: map[string]string{}, keyX: map[string]int{}, keyY: map[string]int{}}
	lines := strings.Split(string(buf), "\n")
	for y, line := range lines {
		if line == "" {
			continue
		}
		puz.width = len(line)
		puz.height++
		for x := 0; x < len(line); x++ {
			if line[x:x+1] == "L" {
				key := genKey(x, y)
				puz.grid[key] = "L"
				puz.keyX[key] = x
				puz.keyY[key] = y
			}
		}
	}

	return puz
}

func (p *Puzzle) lookup(x, y int) string {
	return p.grid[genKey(x, y)]
}

func genKey(x, y int) string {
	return fmt.Sprintf("%v,%v", y, x)
}

func (p *Puzzle) parseKey(key string) (x, y int) {
	return p.keyX[key], p.keyY[key]
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
