// -*- compile-command: "go run main.go ../input.txt"; -*-

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
	groups := readGroups(filename)
	// fmt.Printf("%v\n", strings.Join(groups, "\n"))
	var count int
	for _, group := range groups {
		group = strings.TrimSpace(group)
		// fmt.Printf("%v\n", group)
		runes := map[rune]int{}
		lines := strings.Split(group, "\n")
		for _, line := range lines {
			for _, r := range line {
				runes[r]++
				if runes[r] == len(lines) {
					// fmt.Printf("rune(%c)=%v\n", r, len(lines))
					count++
				}
			}
		}
		// fmt.Printf("len(lines)=%v, count=%v\n", len(lines), count)
	}
	fmt.Printf("Solution: %v\n", count)
}

func readGroups(filename string) []string {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)
	groups := strings.Split(string(buf), "\n\n")
	return groups
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
