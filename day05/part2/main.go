// -*- compile-command: "go run main.go ../input.txt"; -*-

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
	spaces := readLines(filename)
	ids := processSpaces(spaces)
	sort.Sort(sort.IntSlice(ids))

	log.Printf("All ids: %v", ids)

	for i, v := range ids {
		if i > 0 && v-ids[i-1] > 1 {
			log.Printf("My seat: %v", v-1)
			break
		}
	}
}

func readLines(filename string) []string {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)
	lines := strings.Split(string(buf), "\n")
	return lines
}

func processSpaces(lines []string) []int {
	var result []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		v := spaceID(line)
		result = append(result, v)
	}

	return result
}

func spaceID(s string) int {
	rowS := s[0:7]
	colS := s[7:10]
	rowS = strings.ReplaceAll(rowS, "F", "0")
	rowS = strings.ReplaceAll(rowS, "B", "1")
	colS = strings.ReplaceAll(colS, "L", "0")
	colS = strings.ReplaceAll(colS, "R", "1")

	row, err := strconv.ParseInt(rowS, 2, 64)
	check("row ParseInt: %v", err)
	col, err := strconv.ParseInt(colS, 2, 64)
	check("col ParseInt: %v", err)

	return int(row*8 + col)
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
