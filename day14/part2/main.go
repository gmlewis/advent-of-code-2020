// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
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
	var maskFloating int64
	var maskAnd int64
	var maskOr int64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mask = ") {
			bitstring := line[len("mask = "):]
			maskFloating, maskAnd, maskOr = parseMask(bitstring)
			continue
		}
		m := memRE.FindStringSubmatch(line)
		address, err := strconv.ParseInt(m[1], 10, 64)
		check("address: %v", err)
		value, err := strconv.ParseInt(m[2], 10, 64)
		check("value: %v", err)

		bits := []int{}
		for i := 0; i < 36; i++ {
			if maskFloating&(int64(1)<<(35-i)) != 0 {
				bits = append(bits, 35-i)
			}
		}

		total := int(math.Pow(2, float64(len(bits))))
		for i := 0; i < total; i++ {
			newAddress := calcAddress(i, bits, address, maskFloating, maskAnd, maskOr)
			mem[newAddress] = value
		}
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	log.Printf("Solution: %v", sum)
}

func calcAddress(i int, bits []int, address, maskFloating, maskAnd, maskOr int64) int64 {
	var val int64
	for j := 0; j < len(bits); j++ {
		if i&(1<<j) > 0 {
			val |= (1 << bits[j])
		}
	}
	result := (address & maskAnd) | maskOr | val
	// log.Printf("i=%v, bits=%v, val=%v, address=%v, result=%v", i, bits, val, address, result)
	return result
}

func parseMask(s string) (maskFloating, maskAnd, maskOr int64) {
	for i := 0; i < 36; i++ {
		if s[i:i+1] == "X" {
			maskFloating |= int64(1) << (35 - i)
		} else if s[i:i+1] == "1" {
			maskOr |= int64(1) << (35 - i)
		} else {
			maskAnd |= int64(1) << (35 - i)
		}
	}
	return maskFloating, maskAnd, maskOr
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
