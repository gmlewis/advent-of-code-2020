// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")

	memRE = regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)
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

	lines := strings.Split(strings.TrimSpace(string(buf)), "\n")
	mem := map[int64]int64{}
	var maskAnd int64
	var maskOr int64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mask = ") {
			bitstring := line[len("mask = "):]
			maskAnd, maskOr = parseMask(bitstring)
			continue
		}
		m := memRE.FindStringSubmatch(line)
		address, err := strconv.ParseInt(m[1], 10, 64)
		check("address: %v", err)
		value, err := strconv.ParseInt(m[2], 10, 64)
		check("value: %v", err)

		mem[address] = (value & maskAnd) | maskOr
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	log.Printf("Solution: %v", sum)
}

func parseMask(s string) (maskAnd int64, maskOr int64) {
	for i := 0; i < 36; i++ {
		if s[i:i+1] == "X" {
			maskAnd |= int64(1) << (35 - i)
		} else if s[i:i+1] == "1" {
			maskOr |= int64(1) << (35 - i)
		}
	}
	return maskAnd, maskOr
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
