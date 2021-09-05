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

	foodRE = regexp.MustCompile(`^(.*?) \(contains (.*?)\)$`)
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

	players := strings.Split(string(buf), "\n\n")
	p1 := []int{}
	for _, line := range strings.Split(strings.TrimSpace(players[0]), "\n")[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		v, err := strconv.Atoi(line)
		check("v: %v", err)
		p1 = append(p1, v)
	}

	p2 := []int{}
	for _, line := range strings.Split(strings.TrimSpace(players[1]), "\n")[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		v, err := strconv.Atoi(line)
		check("v: %v", err)
		p2 = append(p2, v)
	}

	var c1, c2 int
	for len(p1) > 0 && len(p2) > 0 {
		c1, p1 = p1[0], p1[1:]
		c2, p2 = p2[0], p2[1:]
		if c1 > c2 {
			p1 = append(p1, c1, c2)
		} else if c2 > c1 {
			p2 = append(p2, c2, c1)
		} else {
			log.Fatalf("c1=%v, c2=%v", c1, c2)
		}
	}

	log.Printf("p1=%v, p2=%v", p1, p2)
	winner := p1
	if len(p2) > len(p1) {
		winner = p2
	}

	var score int
	for i, v := range winner {
		value := v * (len(winner) - i)
		score += value
	}
	log.Printf("Solution: %v", score)
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
