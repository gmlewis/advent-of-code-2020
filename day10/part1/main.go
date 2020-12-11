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

	oneDiffs := 1
	threeDiffs := 1

	for i := 1; i < len(ints); i++ {
		switch ints[i] - ints[i-1] {
		case 1:
			oneDiffs++
		case 3:
			threeDiffs++
		}
	}

	log.Printf("oneDiffs=%v, threeDiffs=%v, Solution: %v", oneDiffs, threeDiffs, oneDiffs*threeDiffs)
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
