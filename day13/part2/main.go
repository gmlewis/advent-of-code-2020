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
	ids := map[uint64]uint64{}
	sortedIDs := []uint64{}
	// multipliers := map[uint64]uint64{}

	var max uint64
	var start uint64 = 1
	for i, v := range busIDs {
		if v == "x" {
			continue
		}
		id, err := strconv.ParseUint(v, 10, 64)
		check("id: %v", err)
		ids[id] = uint64(i)
		sortedIDs = append(sortedIDs, id)
		// multipliers[id] = 1
		if id > max {
			max = id
		}
		start *= id
	}
	sort.Slice(sortedIDs, func(a, b int) bool { return sortedIDs[a] > sortedIDs[b] })
	log.Printf("max=%v, start=%v", max, start)
	log.Printf("ids=%#v", ids)
	log.Printf("sortedIDs=%#v", sortedIDs)

	// for k := range ids {
	// 	for k2 := range multipliers {
	// 		if k != k2 {
	// 			multipliers[k2] *= k
	// 		}
	// 	}
	// }
	// log.Printf("multipliers=%#v", multipliers)

	// solve(max, start, ids, sortedIDs)

	offset := ids[max]
	log.Printf("offset=%#v", offset)

	timestamp := start/2 - (start / 2 % max)

	for ; timestamp > start/3; timestamp -= max {
		if isSolution(timestamp, offset, ids) {
			log.Printf("Solution: %v", timestamp-offset)
			return
		}
	}
	log.Fatalf("no solution found")
}

func solve(max, start uint64, ids map[uint64]uint64, sortedIDs []uint64) {
	log.Printf("solve(max=%v, start=%v, ids=%v, sortedIDs=%v", max, start, ids, sortedIDs)
	// var ts uint64
	// var delta uint64
	ch := make(chan uint64, *workers)
	for i, id := range sortedIDs {
		factors := []uint64{}
		for j := uint64(1); j < 10; j++ {
			factors = append(factors, j*id-ids[id])
		}
		log.Printf("id=%v, solution must be at N[%v]*%v-%v %+v",
			id, i, id, ids[id], factors)

		if i == 0 {
			go func(out chan<- uint64, k, d uint64) {
				v := k - d
				for {
					// log.Printf("f(%v): emiting %v", k, v)
					out <- v
					v += k
				}
			}(ch, id, ids[id])
		} else {
			out := make(chan uint64, *workers)
			go func(ch <-chan uint64, out chan<- uint64, k, d uint64) {
				v := k - d
				for next := range ch {
					// log.Printf("f(%v): got %v", k, next)
					for v <= next {
						v += k
					}
					if v == next {
						// log.Printf("f(%v): emiting %v", k, v)
						out <- v
					}
				}
			}(ch, out, id, ids[id])
			ch = out
		}
	}

	v := <-ch
	log.Fatalf("Solution: %v", v)

	// t := uint64(1068781)
	// prod := (t + 4) * (t + 6) * (t + 7) * (t + 1) * t
	// log.Printf("(t+4)(t+6)(t+7)(t+1)(t)=%v %% %v = %v", prod, start, prod%start)
	// log.Printf("(t+4)(t+6)(t+7)(t+1)(t)=%v / %v = %v", prod, start, prod/start)
}

func lcm(a, b uint64) uint64 {
	return a * b / gcd(a, b)
}

func gcd(a, b uint64) uint64 {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func isSolution(timestamp, offset uint64, ids map[uint64]uint64) bool {
	// log.Printf("timestamp=%v", timestamp)
	for k, v := range ids {
		if timestamp == 1068785 {
			log.Printf("(ts=%v - off=%v + v=%v = %v)%%(k=%v) = %v; (%v/%v)=%v",
				timestamp, offset, v, (timestamp - offset + v), k, (timestamp-offset+v)%k,
				(timestamp - offset + v), k, (timestamp-offset+v)/k)
		}
		if (timestamp-offset+v)%k != 0 {
			return false
		}
	}
	return true
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
