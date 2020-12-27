// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
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
	tiles := map[key]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key := parse(line)
		if tiles[key] {
			delete(tiles, key)
		} else {
			tiles[key] = true
		}
	}

	log.Printf("Solution: %v", len(tiles))
}

type key struct {
	x, y int
}

func parse(s string) key {
	var k key
	evenRow := func() bool {
		v := k.y
		if v < 0 {
			v = -v
		}
		return v%2 == 0
	}
	oddRow := func() bool {
		v := k.y
		if v < 0 {
			v = -v
		}
		return v%2 == 1
	}

	for len(s) > 0 {
		switch {
		case strings.HasPrefix(s, "w"):
			k.x--
			s = s[1:]
		case strings.HasPrefix(s, "e"):
			k.x++
			s = s[1:]
		case strings.HasPrefix(s, "se") && evenRow():
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "se") && oddRow():
			k.x++
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "sw") && evenRow():
			k.x--
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "sw") && oddRow():
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "ne") && evenRow():
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "ne") && oddRow():
			k.x++
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "nw") && evenRow():
			k.x--
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "nw") && oddRow():
			k.y++
			s = s[2:]
		default:
			log.Fatalf("parse error: %v", s)
		}
	}
	return k
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
