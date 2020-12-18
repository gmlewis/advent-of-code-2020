// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

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
	workers = flag.Int("w", 1000, "Number of workers")
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
	busIDs := strings.Split(lines[1], ",")
	ids := map[int64]int64{}
	sortedIDs := []int64{}

	var max int64
	var start int64 = 1
	for i, v := range busIDs {
		if v == "x" {
			continue
		}
		id, err := strconv.ParseInt(v, 10, 64)
		check("id: %v", err)
		ids[id] = int64(i)
		sortedIDs = append(sortedIDs, id)
		if id > max {
			max = id
		}
		start *= id
	}
	sort.Slice(sortedIDs, func(a, b int) bool { return sortedIDs[a] > sortedIDs[b] })
	log.Printf("max=%v, start=%v", max, start)
	log.Printf("ids=%#v", ids)
	log.Printf("sortedIDs=%#v", sortedIDs)

	offset := ids[max]
	log.Printf("offset=%#v", offset)

	timestamp := int64(100000000000000) - int64(100000000000000)%max
	log.Printf("starting at timestamp %v", timestamp)

	for {
		skip, ok := isSolution(timestamp, offset, ids)
		if ok {
			log.Printf("Solution: %v", timestamp-offset)
			return
		}
		timestamp += skip
	}
	log.Fatalf("no solution found")
}

func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

func gcd(a, b int64) int64 {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func isSolution(timestamp, offset int64, ids map[int64]int64) (int64, bool) {
	// log.Printf("timestamp=%v", timestamp)
	skip := int64(1)
	solved := true
	for k, v := range ids {
		if timestamp == 1068785 {
			log.Printf("(ts=%v - off=%v + v=%v = %v)%%(k=%v) = %v; (%v/%v)=%v",
				timestamp, offset, v, (timestamp - offset + v), k, (timestamp-offset+v)%k,
				(timestamp - offset + v), k, (timestamp-offset+v)/k)
		}
		if (timestamp-offset+v)%k != 0 {
			solved = false
			continue
		}
		skip *= k
	}
	return skip, solved
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
