// -*- compile-command: "go run main.go -p 5 ../example1.txt && go run main.go -p 25 ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	preamble = flag.Int("p", 25, "Preamble length")
	verbose  = flag.Bool("v", false, "Verbose log messages")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	var ints []int

	lines := readLines(filename)
	logf("lines=%v", lines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		n, err := strconv.Atoi(line)
		check("strconv.Atoi: %v", err)
		ints = append(ints, n)
	}

	for i := *preamble; i < len(ints); i++ {
		if foundSolution(i, ints) {
			break
		}
	}
}

func foundSolution(n int, ints []int) bool {
	v := ints[n]
	seen := map[int]bool{}
	for i := n - *preamble; i < n; i++ {
		d := v - ints[i]
		if seen[d] {
			logf("n=%v, v=%v, seen[%v]", n, v, d)
			return false
		}
		seen[ints[i]] = true
	}
	log.Printf("Part 1 solution: %v", v)

	for i := n - 1; i >= *preamble-2; i-- {
		logf("\nChecking i=%v ...", i)
		var sum int
		min := ints[i]
		max := ints[i]
		for j := 0; j < *preamble-1; j++ {
			if ints[i-j] < min {
				min = ints[i-j]
			}
			if ints[i-j] > max {
				max = ints[i-j]
			}

			sum += ints[i-j]
			logf("i=%v, j=%v, ints[%v]=%v, sum=%v", i, j, i-j, ints[i-j], sum)
			if sum > v {
				break
			}
			if sum == v {
				log.Printf("Part 2 solution: %v + %v = %v", min, max, min+max)
				return true
			}
		}
	}
	return true
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
