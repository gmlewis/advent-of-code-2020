// -*- compile-command: "go run main.go ../example1.txt ../example2.txt ../input.txt"; -*-

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
	ints := map[int]bool{}
	var max int

	lines := readLines(filename)
	logf("lines=%v", lines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		n, err := strconv.Atoi(line)
		check("strconv.Atoi: %v", err)
		ints[n] = true
		if n > max {
			max = n
		}
	}

	cache := map[int]int{}
	count := countPossibilities(0, ints, max, cache)

	log.Printf("Solution: %v", count)
}

func countPossibilities(lastN int, ints map[int]bool, max int, cache map[int]int) int {
	if v, ok := cache[lastN]; ok {
		return v
	}

	if lastN == max {
		// log.Printf("countP(lastN=%v, max=%v): 1 => DONE!!!", lastN, max)
		cache[lastN] = 1
		return 1
	}
	if lastN > max {
		// log.Printf("countP(lastN=%v, max=%v): 0 => not a terminal", lastN, max)
		cache[lastN] = 0
		return 0
	}

	var result int
	for i := 1; i <= 3; i++ {
		if ints[lastN+i] {
			result += countPossibilities(lastN+i, ints, max, cache)
		}
		// log.Printf("countP(lastN=%v, max=%v): i=%v, result=%v", lastN, max, i, result)
	}

	cache[lastN] = result
	return result
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
