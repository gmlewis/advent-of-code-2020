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
	lines := readLines(filename)
	start, err := strconv.Atoi(lines[0])
	check("start: %v", err)
	log.Printf("start=%v", start)
	busIDs := map[int]bool{}
	for _, v := range strings.Split(lines[1], ",") {
		if v == "x" {
			continue
		}
		id, err := strconv.Atoi(v)
		check("id: %v", err)
		busIDs[id] = true
	}
	log.Printf("busIDs=%#v", busIDs)

	for wait := 0; true; wait++ {
		v := start + wait
		for id := range busIDs {
			log.Printf("bus ID=%v - %v %% %v = %v", id, v, id, v%id)
			if v%id == 0 {
				log.Printf("wait=%v, busID=%v, Solution: %v", wait, id, wait*id)
				return
			}
		}
	}
}

func readLines(filename string) []string {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)
	return strings.Split(string(buf), "\n")
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
