// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
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

	for {
		logf("grid:\n%v", puz)
		newGrid := puz.iterate()
		if reflect.DeepEqual(newGrid, puz.grid) {
			break
		}
		puz.grid = newGrid
		logf("newGrid:\n%v", puz)
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

func (p *Puzzle) iterate() map[string]string {
	r := map[string]string{}
	for k, v := range p.grid {
		adj := p.countAdjacentOccupied(k)
		// log.Printf("k=%v, v=%v, adj=%v", k, v, adj)
		if v == "L" && adj == 0 {
			r[k] = "#"
		} else if v == "#" && adj >= 4 {
			r[k] = "L"
		} else {
			r[k] = v
		}
	}
	return r
}

func (p *Puzzle) countAdjacentOccupied(k string) int {
	x, y := parseKey(k)
	var adj int

	t := func(x, y int) {
		v := p.lookup(x, y)
		if v == "#" {
			adj++
		}
	}

	t(x-1, y-1)
	t(x, y-1)
	t(x+1, y-1)
	t(x-1, y)

	t(x+1, y)
	t(x-1, y+1)
	t(x, y+1)
	t(x+1, y+1)

	return adj
}

type Puzzle struct {
	width  int
	height int
	grid   map[string]string
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

	puz := &Puzzle{grid: map[string]string{}}
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

func parseKey(key string) (x, y int) {
	parts := strings.Split(key, ",")
	var err error
	x, err = strconv.Atoi(parts[1])
	check("key=%q parseKey(%q).x: %v", key, parts[1], err)
	y, err = strconv.Atoi(parts[0])
	check("key=%q parseKey(%q).y: %v", key, parts[0], err)
	return x, y
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
