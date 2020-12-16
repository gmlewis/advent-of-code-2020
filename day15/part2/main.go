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

	type record struct {
		count      int
		mostRecent int
		beforeThat int
	}
	// whenSpoken lists the last times this value was spoken
	whenSpoken := map[int]*record{}

	speakIt := func(n, turn int) {
		lastSpoken = n
		if v, ok := whenSpoken[n]; ok {
			v.count++
			v.beforeThat = v.mostRecent
			v.mostRecent = turn
		} else {
			whenSpoken[n] = &record{count: 1, mostRecent: turn}
		}
	}

	for _, part := range parts {
		v, err := strconv.Atoi(part)
		check("v: %v", err)
		turn++
		speakIt(v, turn)
	}
	// log.Printf("whenSpoken: %#v", whenSpoken)

	for {
		numLastSpoken := whenSpoken[lastSpoken].count
		turn++
		var speak int
		if numLastSpoken > 1 {
			speak = whenSpoken[lastSpoken].mostRecent - whenSpoken[lastSpoken].beforeThat
		}
		if turn == 30000000 {
			log.Printf("Turn %v: %v", turn, speak)
			break
		}
		speakIt(speak, turn)
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
