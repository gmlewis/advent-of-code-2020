// -*- compile-command: "go run main.go ../example1.txt ../example2.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"sort"
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

	sort.Sort(sort.IntSlice(ints))

	for f1 := 0; f1 <= 9; f1++ {
		for f2 := 0; f2 <= 9; f2++ {
			s1 := findPossibilities(2, ints, f1, f2)
			if s1 <= 19208 {
				log.Printf("OFFSET Solution(%v, %v): %v", f1, f2, s1)
			}
			s2 := findPossibilities(0, ints, f1, f2)
			if s2 <= 19208 {
				log.Printf("Solution(%v, %v): %v DIFF=%v", f1, f2, s2, 19208-s2)
			}
		}
	}
}

func findPossibilities(n int, ints []int, f1, f2 int) int {
	count := 1
	if n < 0 || n >= len(ints) {
		return count
	}
	// log.Printf("ENTER: findPossibilities(ints[%v]=%v)", n, ints[n])

	var diffs string
	for i := 1; n+i < len(ints); i++ {
		if ints[n+i]-ints[n+i-1] == 1 {
			diffs += "1"
		} else {
			break
		}
	}

	switch diffs {
	case "":
		// log.Printf("findPossibilities(ints[%v]=%v) = %q = 1", n, ints[n], diffs)
		return findPossibilities(n+1+len(diffs), ints, f1, f2)
	case "1":
		// log.Printf("findPossibilities(ints[%v]=%v) = %q = 1", n, ints[n], diffs)
		return findPossibilities(n+1+len(diffs), ints, f1, f2)
	case "11":
		// log.Printf("findPossibilities(ints[%v]=%v) = %q = 2", n, ints[n], diffs)
		return 2 * findPossibilities(n+1+len(diffs), ints, f1, f2)
	case "111":
		// log.Printf("findPossibilities(ints[%v]=%v) = %q = 4", n, ints[n], diffs)
		return 4 * findPossibilities(n+1+len(diffs), ints, f1, f2)
	case "1111":
		// log.Printf("BEFORE findPossibilities(ints[%v]=%v) = %q", n, ints[n], diffs)
		p1 := findPossibilities(n+2, ints, f1, f2)
		p2 := findPossibilities(n+1+len(diffs), ints, f1, f2)
		// log.Printf("AFTER findPossibilities(ints[%v]=%v) = %q: %v*%v+%v*%v", n, ints[n], diffs, f2, p2, f1, p1)
		return f2*p2 + f1*p1
	default:
		log.Fatalf("UNHANDLED: findPossibilities(ints[%v]=%v) = %v", n, ints[n], diffs)
	}
	return 0
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
