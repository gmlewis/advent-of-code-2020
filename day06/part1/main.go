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
		group = strings.ReplaceAll(group, "\n", "")
		// fmt.Printf("%v\n", group)
		runes := map[rune]int{}
		for _, r := range group {
			runes[r]++
		}
		// fmt.Printf("%v\n", len(runes))
		count += len(runes)
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
