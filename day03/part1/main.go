package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	down    = flag.Int("down", 1, "Number of cells to move down")
	right   = flag.Int("right", 3, "Number of cells to move right")
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

	var posX, posY, count int
	for y := 0; y < puz.height; y++ {
		if puz.lookup(posX, posY) {
			count++
		}
		posX += *right
		posY += *down
	}

	log.Printf("Found %v trees", count)
}

type Puzzle struct {
	width  int
	height int
	grid   map[string]bool
}

func readPuzzle(filename string) *Puzzle {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	puz := &Puzzle{grid: map[string]bool{}}
	lines := strings.Split(string(buf), "\n")
	for y, line := range lines {
		if line == "" {
			continue
		}
		puz.width = len(line)
		puz.height++
		for x := 0; x < len(line); x++ {
			if line[x:x+1] == "#" {
				key := genKey(x, y)
				puz.grid[key] = true
			}
		}
	}

	return puz
}

func (p *Puzzle) lookup(x, y int) bool {
	u := x % p.width
	// v := y % p.height
	return p.grid[genKey(u, y)]
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
