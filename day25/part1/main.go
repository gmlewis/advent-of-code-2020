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
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	lines := strings.Split(string(buf), "\n")
	cpk, err := strconv.Atoi(lines[0])
	dpk, err := strconv.Atoi(lines[1])
	log.Printf("cpk=%v, dpk=%v", cpk, dpk)
	cls := loopSize(7, cpk)
	dls := loopSize(7, dpk)
	log.Printf("cls=%v, dls=%v", cls, dls)
	encKey := encryptionKey(cpk, dls)

	log.Printf("Solution: %v", encKey)
}

func encryptionKey(subNum, loopSize int) int {
	result := 1
	for i := 0; i < loopSize; i++ {
		result *= subNum
		result = result % 20201227
	}
	return result
}

func loopSize(subNum int, pk int) int {
	result := 1
	lastVal := 1
	for {
		v := lastVal
		v *= subNum
		v = v % 20201227
		if v == pk {
			return result
		}
		lastVal = v
		result++
	}
	return 0
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
