// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

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
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	parts := strings.Split(strings.TrimSpace(string(buf)), ",")
	var turn int
	var lastSpoken int
	// whenSpoken lists the last times this value was spoken
	whenSpoken := map[int][]int{}
	for _, part := range parts {
		v, err := strconv.Atoi(part)
		check("v: %v", err)
		turn++
		lastSpoken = v
		whenSpoken[lastSpoken] = append(whenSpoken[lastSpoken], turn)
	}

	for {
		numLastSpoken := len(whenSpoken[lastSpoken])
		turn++
		var speak int
		if numLastSpoken > 1 {
			speak = whenSpoken[lastSpoken][numLastSpoken-1] - whenSpoken[lastSpoken][numLastSpoken-2]
		}
		if turn == 2020 {
			log.Printf("Turn %v: %v", turn, speak)
			break
		}
		lastSpoken = speak
		whenSpoken[lastSpoken] = append(whenSpoken[lastSpoken], turn)
	}

	// log.Printf("Solution: %v", solution)
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
