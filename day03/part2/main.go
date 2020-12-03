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
	verbose = flag.Bool("v", false, "Verbose log messages")

	rights = []int{1, 3, 5, 7, 1}
	downs  = []int{1, 1, 1, 1, 2}
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

	result := 1
	for i := 0; i < len(rights); i++ {
		count := puz.countTrees(rights[i], downs[i])
		result *= count
	}

	log.Printf("Solution: %v", result)
}

func (p *Puzzle) countTrees(right, down int) int {
	var posX, posY, count int
	for y := 0; y < p.height; y++ {
		if p.lookup(posX, posY) {
			count++
		}
		posX += right
		posY += down
	}

	log.Printf("(%v,%v): Found %v trees", right, down, count)
	return count
}

type Puzzle struct {
	width  int
	height int
	grid   map[string]bool
}

func readPuzzle(filename string) *Puzzle {
	logf("Processing %v ...", filename)
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
