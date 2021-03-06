package main

import (
	"flag"
	"fmt"
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

	vals := map[int]bool{}
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		v, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		if find(v, vals) {
			break
		}
		vals[2020-v] = true
	}
}

func find(v int, vals map[int]bool) bool {
	for d2 := range vals {
		v2 := 2020 - d2
		if vals[v+v2] {
			fmt.Printf("%v + %v + %v = 2020\n%v * %v * %v = %v\n", v, v2, 2020-v-v2, v, v2, 2020-v-v2, (v * v2 * (2020 - v - v2)))
			return true
		}
	}
	return false
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
