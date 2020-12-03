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
	lineRE  = regexp.MustCompile(`^(\d+)-(\d+) ([a-z]): (.*)$`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}
}

func process(filename string) {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	var count int
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}

		m := lineRE.FindStringSubmatch(line)
		if len(m) != 5 {
			log.Fatalf("unexpected line: %v", line)
		}

		start, err := strconv.Atoi(m[1])
		check("start strconv.Atoi: %v", err)
		end, err := strconv.Atoi(m[2])
		check("end strconv.Atoi: %v", err)
		letter := m[3]
		passwd := m[4]

		if valid(start, end, letter, passwd) {
			count++
		}
	}

	log.Printf("%v valid passwords", count)
}

func valid(start, end int, letter, passwd string) bool {
	if end > len(passwd) {
		log.Fatalf("end=%v >= len(passwd)=%v: %v", end, len(passwd), passwd)
	}

	first := passwd[start-1:start] == letter
	second := passwd[end-1:end] == letter

	// if first != second {
	// 	log.Printf("%v-%v %v: %v", start, end, letter, passwd)
	// }

	return first != second
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
